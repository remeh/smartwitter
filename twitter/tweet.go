package twitter

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/remeh/uuid"
)

var (
	tweetIdSpace     = uuid.UUID("139866bf-4932-4039-91f1-c8c4a5994837")
	tweetUserIdSpace = uuid.UUID("4ab3e860-8381-4c83-acce-dbb79df2fbc5")
)

type Tweets []Tweet

type Tweet struct {
	// Please use Uid() to gets the UUID of this tweet.
	uid uuid.UUID
	// Time of the entry in the database.
	CreationTime time.Time
	// Id of the tweet on Twitter.
	TwitterId int64
	// Twitter profile creation time.
	TwitterCreationTime time.Time
	Text                string
	UserUid             uuid.UUID
}

type TweetUser struct {
	// Please use Uid() to gets the UUID of this tweet.
	uid uuid.UUID
	// Time of the entry in the database.
	CreationTime time.Time
	// Id of the user on Twitter.
	TwitterId int64
	// Twitter profile creation time.
	TwitterCreationTime time.Time
	Description         string
	ScreenName          string
	Name                string
	TimeZone            string
}

func (t Tweet) Uid() uuid.UUID {
	return uuid.NewSHA1(tweetIdSpace, []byte(fmt.Sprintf("%d", t.TwitterId)))
}

func (t TweetUser) Uid() uuid.UUID {
	return uuid.NewSHA1(tweetUserIdSpace, []byte(fmt.Sprintf("%d", t.TwitterId)))
}

func (t *Tweet) Upsert(db *sql.DB) error {

	if user, err := t.User(db); err != nil {
		return err
	} else {
		// upsert the user
		if err = user.Upsert(db); err != nil {
			return err
		}
	}

	if _, err := db.Exec(`
		INSERT INTO "tweet" ("uid", "creation_time", "twitter_id", "twitter_creation_time", "text", "user_uid")
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT ("uid") DO UPDATE SET
			"creation_time" = $2,
			"twitter_id" = $3,
			"twitter_creation_time" = $4,
			"text" = $5,
			"user_uid" = $6
	`, t.Uid(), t.TwitterId, t.TwitterCreationTime, t.Text, t.UserUid); err != nil {
		return err
	}

	return nil
}

func (t *Tweet) User(db *sql.DB) (*TweetUser, error) {
	return FindTweetUser(db, t.UserUid)
}

func (tu *TweetUser) Upsert(db *sql.DB) error {

	if _, err := db.Exec(`
		INSERT INTO "tweet_user" ("uid", "creation_time", "twitter_id", "twitter_creation_time", "description", "screen_name", "name", "timezone")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT ("uid") DO UPDATE SET
			"creation_time" = $2,
			"twitter_id" = $3,
			"twitter_creation_time" = $4,
			"description" = $5,
			"screen_name" = $6,
			"name" = $7,
			"timezone" = $8
	`, tu.Uid(), tu.CreationTime, tu.TwitterId, tu.TwitterCreationTime, tu.Description, tu.ScreenName, tu.Name, tu.TimeZone); err != nil {
		return err
	}

	return nil
}

// FindTweetUser finds th
func FindTweetUser(db *sql.DB, id uuid.UUID) (*TweetUser, error) {
	if db == nil || id == nil {
		return nil, fmt.Errorf("nil db or nil id")
	}

	rv := &TweetUser{uid: id}

	if err := db.QueryRow(`
		SELECT "creation_time", "twitter_id", "twitter_creation_time", "description", "screen_name", "name", "timezone" FROM "tweet_user"
		WHERE "uid" = $1
		LIMIT 1
	`).Scan(&rv.CreationTime, &rv.TwitterId, &rv.TwitterCreationTime, &rv.Description, &rv.ScreenName, &rv.Name, &rv.TimeZone); err != nil {
		return nil, err
	}

	return rv, nil
}

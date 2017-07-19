package twitter

import (
	"database/sql"
	"fmt"

	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/uuid"
)

type twitterUserDAO struct {
	DB *sql.DB
}

// ----------------------

var tuDao *twitterUserDAO

func TwitterUserDAO() *twitterUserDAO {
	if tuDao != nil {
		return tuDao
	}

	tuDao = &twitterUserDAO{
		DB: storage.DB(),
	}

	if err := tuDao.InitStmt(); err != nil {
		log.Error("Can't prepare TwitterUserDAO")
		panic(err)
	}

	return tuDao
}

func (d *twitterUserDAO) InitStmt() error {
	var err error
	return err
}

func (d *twitterUserDAO) Upsert(tu *TwitterUser) error {
	if tu == nil {
		return fmt.Errorf("tu == nil")
	}

	if _, err := d.DB.Exec(`
		INSERT INTO "twitter_user" ("uid", "creation_time", "last_update", "twitter_id", "description", "screen_name", "name", "timezone", "avatar", "utc_offset", "followers_count")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT ("uid") DO UPDATE SET
			"last_update" = $3,
			"twitter_id" = $4,
			"description" = $5,
			"screen_name" = $6,
			"name" = $7,
			"timezone" = $8,
			"avatar" = $9,
			"utc_offset" = $10,
			"followers_count" = $11
	`, tu.Uid(), tu.CreationTime, tu.LastUpdate, tu.TwitterId, tu.Description, tu.ScreenName, tu.Name, tu.TimeZone, tu.Avatar, tu.UtcOffset, tu.FollowersCount); err != nil {
		return err
	}

	return nil
}

func (d *twitterUserDAO) Find(id uuid.UUID) (*TwitterUser, error) {
	if id == nil {
		return nil, fmt.Errorf("nil db or nil id")
	}

	rv := &TwitterUser{uid: id}
	if err := d.DB.QueryRow(`
		SELECT "creation_time", "last_update", "twitter_id", "description", "screen_name", "name", "timezone", "avatar", "utc_offset", "followers_count" FROM "twitter_user"
		WHERE "uid" = $1
		LIMIT 1
	`, id).Scan(&rv.CreationTime, &rv.LastUpdate, &rv.TwitterId, &rv.Description, &rv.ScreenName, &rv.Name, &rv.TimeZone, &tu.Avatar, &rv.UtcOffset, &rv.FollowersCount); err != nil {
		return nil, err
	}
	return rv, nil
}

package twitter

import (
	"fmt"
	"time"

	"github.com/remeh/uuid"
)

var (
	tweetIdSpace = uuid.UUID("139866bf-4932-4039-91f1-c8c4a5994837")
)

type Tweets []Tweet

type Tweet struct {
	// Please use Uid() to gets the UUID of this tweet.
	uid uuid.UUID `json:"-"`
	// Time of the entry in the database.
	CreationTime time.Time `json:"creation_time"`
	LastUpdate   time.Time `json:"last_update"`
	// Id of the tweet on Twitter.
	TwitterId int64 `json:"-"`
	// Twitter profile creation time.
	TwitterCreationTime time.Time `json:"twitter_creation_time"`
	RetweetCount        int       `json:"retweet_count"`
	FavoriteCount       int       `json:"favorite_count"`
	Text                string    `json:"text"`
	Lang                string    `json:"-"`
	Link                string    `json:"link"`
	UserUid             uuid.UUID `json:"-"`
	// keywords having found this tweet
	Keywords []string `json:"-"`
}

func (t Tweet) Uid() uuid.UUID {
	if t.uid == nil && t.TwitterId >= 0 {
		t.uid = uuid.NewSHA1(tweetIdSpace, []byte(fmt.Sprintf("%d", t.TwitterId)))
	}
	return t.uid
}

func NewTweet(twitterId, twitterUserId int64) *Tweet {
	return &Tweet{
		TwitterId: twitterId,
		UserUid:   NewTwitterUser(twitterUserId).Uid(),
	}
}

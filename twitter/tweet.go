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
	uid uuid.UUID
	// Time of the entry in the database.
	CreationTime time.Time
	LastUpdate   time.Time
	// Id of the tweet on Twitter.
	TwitterId int64
	// Twitter profile creation time.
	TwitterCreationTime time.Time
	RetweetCount        int
	FavoriteCount       int
	Text                string
	UserUid             uuid.UUID
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

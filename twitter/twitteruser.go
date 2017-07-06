package twitter

import (
	"fmt"
	"time"

	"github.com/remeh/uuid"
)

var (
	twitterUserIdSpace = uuid.UUID("4ab3e860-8381-4c83-acce-dbb79df2fbc5")
)

type TwitterUsers []TwitterUser

type TwitterUser struct {
	// Please use Uid() to gets the UUID of this tweet.
	uid uuid.UUID
	// Time of the entry in the database.
	CreationTime time.Time
	LastUpdate   time.Time
	// Id of the user on Twitter.
	TwitterId      int64
	Description    string
	ScreenName     string
	Name           string
	TimeZone       string
	UtcOffset      int
	FollowersCount int
}

func (t TwitterUser) Uid() uuid.UUID {
	if t.uid == nil && t.TwitterId >= 0 {
		t.uid = uuid.NewSHA1(twitterUserIdSpace, []byte(fmt.Sprintf("%d", t.TwitterId)))
	}
	return t.uid
}

func NewTwitterUser(twitterId int64) *TwitterUser {
	return &TwitterUser{
		TwitterId: twitterId,
	}
}

package twitter

import (
	"fmt"
	"time"

	"github.com/remeh/uuid"
)

var (
	twitterUserIdSpace = uuid.UUID("4ab3e860-8381-4c83-acce-dbb79df2fbc5")
)

type TwitterUser struct {
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

func (t TwitterUser) Uid() uuid.UUID {
	return uuid.NewSHA1(twitterUserIdSpace, []byte(fmt.Sprintf("%d", t.TwitterId)))
}

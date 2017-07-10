package twitter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/remeh/anaconda"
	"github.com/remeh/smartwitter/log"
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
	TwitterId string `json:"-"`
	// Twitter profile creation time.
	TwitterCreationTime time.Time `json:"twitter_creation_time"`
	RetweetCount        int       `json:"retweet_count"`
	FavoriteCount       int       `json:"favorite_count"`
	Text                string    `json:"text"`
	Lang                string    `json:"-"`
	Link                string    `json:"link"`
	TwitterUserUid      uuid.UUID `json:"-"`
	// keywords having found this tweet
	Keywords []string `json:"-"`
}

func (t Tweet) Uid() uuid.UUID {
	if t.uid == nil && t.TwitterId != "" {
		t.uid = uuid.NewSHA1(tweetIdSpace, []byte(fmt.Sprintf("%s", t.TwitterId)))
	}
	return t.uid
}

func (t Tweet) User() (*TwitterUser, error) {
	tu, err := TwitterUserDAO().Find(t.TwitterUserUid)
	if err != nil {
		return nil, err
	}
	if tu == nil {
		return nil, fmt.Errorf("can't find user: %v", t.TwitterUserUid)
	}
	return tu, nil
}

func NewTweet(twitterId, twitterUserId string) *Tweet {
	return &Tweet{
		TwitterId:      twitterId,
		TwitterUserUid: NewTwitterUser(twitterUserId).Uid(),
	}
}

func TweetFromTweet(t anaconda.Tweet, now time.Time, keywords []string) *Tweet {
	var err error
	tid := strconv.FormatInt(t.Id, 10)
	rv := NewTweet(tid, strconv.FormatInt(t.User.Id, 10))
	if rv.TwitterCreationTime, err = t.CreatedAtTime(); err != nil {
		log.Warning("getTweets: getting tweet creation time:", err)
	}
	rv.CreationTime = now
	rv.LastUpdate = now
	rv.Text = t.FullText
	rv.RetweetCount = t.RetweetCount
	rv.FavoriteCount = t.FavoriteCount
	rv.Lang = t.Lang
	rv.Keywords = keywords
	rv.Link = fmt.Sprintf("https://twitter.com/%s/status/%s", t.User.ScreenName, tid)
	return rv
}

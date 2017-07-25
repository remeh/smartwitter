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

// Tweet Entities
// ----------------------

type TweetEntities []TweetEntity

type TweetEntity struct {
	Type          TweetEntityType `json:"type,omitempty"`
	DisplayUrl    string          `json:"display_url,omitempty"`
	Url           string          `json:"url,omitempty"`
	Indices       []int           `json:"indices,omitempty"`
	ScreenName    string          `json:"screen_name,omitempty"`
	UserName      string          `json:"user_name,omitempty"`
	UserTwitterId string          `json:"user_twitter_id,omitempty"`
	Hashtag       string          `json:"hashtag,omitempty"`
}

type TweetEntityType string

const (
	Media       TweetEntityType = "media"
	Hashtag     TweetEntityType = "hashtag"
	Url         TweetEntityType = "url"
	UserMention TweetEntityType = "user_mention"
)

func ToTweetEntities(types, displayUrls, urls []string, idxStarts, idxEnds []int64, screenNames, userNames, userTwitterIds, hashtags []string) TweetEntities {
	if len(types) != len(displayUrls) ||
		len(types) != len(urls) ||
		len(types) != len(idxStarts) ||
		len(types) != len(idxEnds) ||
		len(types) != len(screenNames) ||
		len(types) != len(userNames) ||
		len(types) != len(userTwitterIds) ||
		len(types) != len(hashtags) {
		log.Error("ToTweetEntities: one of the slices isn't of the good size.")
		return make(TweetEntities, 0)
	}

	rv := make(TweetEntities, len(types))

	for i := range types {
		rv[i] = TweetEntity{
			Type:          TweetEntityType(types[i]),
			DisplayUrl:    displayUrls[i],
			Url:           urls[i],
			Indices:       []int{int(idxStarts[i]), int(idxEnds[i])},
			ScreenName:    screenNames[i],
			UserName:      userNames[i],
			UserTwitterId: userTwitterIds[i],
			Hashtag:       hashtags[i],
		}
	}

	return rv
}

func (t TweetEntities) Types() []string {
	if t == nil {
		return nil
	}
	rv := make([]string, 0)
	for _, ent := range t {
		rv = append(rv, string(ent.Type))
	}
	return rv
}

func (t TweetEntities) DisplayUrls() []string {
	if t == nil {
		return nil
	}
	rv := make([]string, 0)
	for _, ent := range t {
		rv = append(rv, ent.DisplayUrl)
	}
	return rv
}

func (t TweetEntities) Urls() []string {
	if t == nil {
		return nil
	}
	rv := make([]string, 0)
	for _, ent := range t {
		rv = append(rv, ent.Url)
	}
	return rv
}

func (t TweetEntities) Indices() ([]int, []int) {
	if t == nil {
		return nil, nil
	}
	rvs, rve := make([]int, 0), make([]int, 0)
	for _, ent := range t {
		if len(ent.Indices) == 2 {
			rvs = append(rvs, ent.Indices[0])
			rve = append(rve, ent.Indices[1])
		}
	}
	return rvs, rve
}

func (t TweetEntities) ScreenNames() []string {
	if t == nil {
		return nil
	}
	rv := make([]string, 0)
	for _, ent := range t {
		rv = append(rv, ent.ScreenName)
	}
	return rv
}

func (t TweetEntities) UserNames() []string {
	if t == nil {
		return nil
	}
	rv := make([]string, 0)
	for _, ent := range t {
		rv = append(rv, ent.UserName)
	}
	return rv
}

func (t TweetEntities) UserTwitterIds() []string {
	if t == nil {
		return nil
	}
	rv := make([]string, 0)
	for _, ent := range t {
		rv = append(rv, ent.UserTwitterId)
	}
	return rv
}

func (t TweetEntities) Hashtags() []string {
	if t == nil {
		return nil
	}
	rv := make([]string, 0)
	for _, ent := range t {
		rv = append(rv, ent.Hashtag)
	}
	return rv
}

// Tweets
// ----------------------

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

	Entities TweetEntities `json:"entities"`
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

	// Entities
	// ----------------------

	ents := make(TweetEntities, 0)

	for _, h := range t.Entities.Hashtags {
		ents = append(ents, TweetEntity{
			Type:    Hashtag,
			Indices: h.Indices[:],
		})
	}

	for _, m := range t.Entities.Media {
		ents = append(ents, TweetEntity{
			Type:       Media,
			DisplayUrl: m.Display_url,
			Url:        m.Media_url_https,
			Indices:    m.Indices[:],
		})
	}

	for _, um := range t.Entities.User_mentions {
		ents = append(ents, TweetEntity{
			Type:          UserMention,
			UserName:      um.Name,
			ScreenName:    um.Screen_name,
			Indices:       um.Indices[:],
			UserTwitterId: um.Id_str,
		})
	}

	for _, u := range t.Entities.Urls {
		ents = append(ents, TweetEntity{
			Type:       Url,
			DisplayUrl: u.Display_url,
			Url:        u.Expanded_url,
			Indices:    u.Indices[:],
		})
	}

	rv.Entities = ents

	return rv
}

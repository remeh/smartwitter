package suggest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/suggest"
	"github.com/remeh/smartwitter/twitter"
)

type Tweets struct{}

type byKeywords []byKeywordsInfo

type byKeywordsInfo struct {
	Uid                  string                `json:"uid"`
	TweetId              string                `json:"tweet_id"`
	Link                 string                `json:"link"`
	Name                 string                `json:"name"`
	ScreenName           string                `json:"screen_name"`
	Avatar               string                `json:"avatar"`
	Time                 time.Time             `json:"time"`
	Text                 string                `json:"text"`
	RetweetCount         int                   `json:"retweet_count"`
	LikeCount            int                   `json:"like_count"`
	TwitterUserFollowers int                   `json:"twitter_user_followers"`
	Ignored              bool                  `json:"ignored"`
	Liked                bool                  `json:"liked"`
	Retweeted            bool                  `json:"retweeted"`
	Entities             twitter.TweetEntities `json:"entities"`
}

func (c Tweets) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// user
	// ----------------------

	var user *account.User
	var err error

	if user, err = api.GetUser(r); err != nil {
		api.RenderErrJson(w, err)
		return
	} else if user == nil {
		api.RenderForbiddenJson(w)
		return
	}

	// parse form
	// ----------------------

	r.ParseForm()
	pp := r.Form.Get("p")
	if len(pp) == 0 {
		api.RenderBadParameter(w, "p")
		return
	}
	p, err := strconv.Atoi(pp)
	if err != nil || p < 0 {
		api.RenderBadParameter(w, "p")
		return
	}

	// get the user keywords for this position
	// ----------------------

	kwords, err := suggest.GetKeywords(user.Uid)
	if len(kwords)-1 < p {
		api.RenderBadParameter(w, "p")
		return
	}

	// get the suggestion
	// ----------------------

	since := time.Hour * 4

	tweets, tdas, err := suggest.SuggestByKeywords(user, kwords[p].Keywords, since, 5)

	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// render the response
	// ----------------------

	data := make(byKeywords, 0)

	for _, t := range tweets {
		tu, err := t.User()
		if err != nil {
			api.RenderErrJson(w, err)
			return
		}

		tda, _ := tdas.Get(t.TwitterId)

		data = append(data, byKeywordsInfo{
			Uid:                  t.Uid().String(),
			TweetId:              t.TwitterId,
			Link:                 t.Link,
			Name:                 tu.Name,
			ScreenName:           tu.ScreenName,
			Avatar:               tu.Avatar,
			Time:                 t.TwitterCreationTime,
			Text:                 t.Text,
			RetweetCount:         t.RetweetCount,
			LikeCount:            t.FavoriteCount,
			TwitterUserFollowers: tu.FollowersCount,
			Ignored:              tda.Ignored,
			Liked:                tda.Liked,
			Retweeted:            tda.Retweeted,
			Entities:             t.Entities,
		})
	}

	api.RenderJson(w, 200, data)
}

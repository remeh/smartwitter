package suggest

import (
	"net/http"
	"time"

	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/suggest"
)

type ByKeywords struct {
}

type byKeywords []byKeywordsInfo

type byKeywordsInfo struct {
	Uid                  string `json:"uid"`
	TweetId              string `json:"tweet_id"`
	Link                 string `json:"link"`
	ScreenName           string `json:"screen_name"`
	Text                 string `json:"text"`
	RetweetCount         int    `json:"retweet_count"`
	FavoriteCount        int    `json:"favorite_count"`
	TwitterUserFollowers int    `json:"twitter_user_followers"`
	Ignored              bool   `json:"liked"`
	Liked                bool   `json:"liked"`
	Retweeted            bool   `json:"retweeted"`
}

func (c ByKeywords) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// parse form
	// ----------------------

	r.ParseForm()
	keywords := api.Escape(r.Form["k"])

	// get the suggestion
	// ----------------------

	since := time.Hour * 2

	tweets, tdas, err := suggest.SuggestByKeywords(keywords, since, 5)

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

		tda := tdas.Get(t.TwitterId)

		data = append(data, byKeywordsInfo{
			Uid:                  t.Uid().String(),
			TweetId:              t.TwitterId,
			Link:                 t.Link,
			ScreenName:           tu.ScreenName,
			Text:                 t.Text,
			RetweetCount:         t.RetweetCount,
			FavoriteCount:        t.FavoriteCount,
			TwitterUserFollowers: tu.FollowersCount,
			Ignored:              tda.Ignored,
			Liked:                tda.Liked,
			Retweeted:            tda.Retweeted,
		})
	}

	api.RenderJson(w, 200, data)
}

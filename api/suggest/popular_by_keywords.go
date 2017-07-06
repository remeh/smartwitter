package suggest

import (
	"net/http"

	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/suggest"
)

type ByKeywords struct {
}

type byKeywords []byKeywordsInfo

type byKeywordsInfo struct {
	Link                 string `json:"link"`
	ScreenName           string `json:"screen_name"`
	Text                 string `json:"text"`
	RetweetCount         int    `json:"retweet_count"`
	FavoriteCount        int    `json:"favorite_count"`
	TwitterUserFollowers int    `json:"twitter_user_followers"`
}

func (c ByKeywords) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// parse form
	// ----------------------

	r.ParseForm()
	keywords := r.Form["k"]

	// TODO(remy): trim all keywords

	// get the suggestion
	// ----------------------

	tweets, err := suggest.SuggestByKeywords(keywords)

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

		data = append(data, byKeywordsInfo{
			Link:                 t.Link,
			ScreenName:           tu.ScreenName,
			Text:                 t.Text,
			RetweetCount:         t.RetweetCount,
			FavoriteCount:        t.FavoriteCount,
			TwitterUserFollowers: tu.FollowersCount,
		})
	}

	api.RenderJson(w, 200, data)
}

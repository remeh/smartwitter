package action

import (
	"net/http"
	"strconv"
	"time"

	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/twitter"
)

type Like struct {
}

func (c Like) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// parse form
	// ----------------------

	r.ParseForm()
	ptid := api.EscapeOne(r.Form.Get("tid"))
	au := api.EscapeOne(r.Form.Get("au")) // auto-unlike

	tid, err := strconv.ParseInt(ptid, 10, 64)
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// ----------------------

	// fav on Twitter
	if _, err := twitter.GetApi().Favorite(tid); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// plan an action to automatically unlike
	if au == "true" {
		unfav := twitter.UnFavorite{
			TweetId: tid,
		}
		unfav.CreationTime = time.Now()
		unfav.ExecutionTime = time.Now().Add(time.Hour * 144) // 6 days

		if err := unfav.Store(); err != nil {
			api.RenderErrJson(w, err)
			return
		}
	}

	api.RenderOk(w)
}

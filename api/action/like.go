package action

import (
	"net/http"
	"strconv"
	"time"

	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/twitter"
	"github.com/remeh/uuid"
)

type Like struct {
}

func (c Like) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// parse form
	// ----------------------

	r.ParseForm()
	ptid := api.EscapeOne(r.Form.Get("tid"))
	au := api.EscapeOne(r.Form.Get("au")) // auto-unlike

	if len(ptid) == 0 {
		api.RenderBadParameter(w, "tid")
		return
	}

	tid, err := strconv.ParseInt(ptid, 10, 64)
	if err != nil {
		api.RenderBadParameter(w, "tid")
		return
	}

	if len(au) == 0 && au != "true" && au != "false" {
		api.RenderBadParameter(w, "au")
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
		unlike := twitter.UnLike{
			Uid:     uuid.New(),
			TweetId: ptid,
		}
		unlike.CreationTime = time.Now()
		unlike.ExecutionTime = time.Now().Add(time.Hour * 144) // 6 days

		if err := unlike.Store(); err != nil {
			api.RenderErrJson(w, err)
			return
		}
	}

	api.RenderOk(w)
}

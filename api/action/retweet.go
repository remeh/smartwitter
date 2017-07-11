package action

import (
	"net/http"
	"strconv"
	"time"

	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/twitter"
	"github.com/remeh/uuid"
)

type Retweet struct {
}

func (c Retweet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// parse form
	// ----------------------

	r.ParseForm()
	ptid := api.EscapeOne(r.Form.Get("tid"))
	au := api.EscapeOne(r.Form.Get("au")) // auto-unretweet

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
	if _, err := twitter.GetApi().Retweet(tid, true); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// plan an action to automatically unrt
	if au == "true" {
		unrt := twitter.UnRetweet{
			Uid:     uuid.New(),
			TweetId: ptid,
		}
		unrt.CreationTime = time.Now()
		unrt.ExecutionTime = time.Now().Add(time.Hour * 144) // 6 days

		if err := unrt.Store(); err != nil {
			api.RenderErrJson(w, err)
			return
		}
	}

	api.RenderOk(w)
}

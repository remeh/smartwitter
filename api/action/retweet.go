package action

import (
	"net/http"
	"strconv"
	"time"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/twitter"
	"github.com/remeh/uuid"
)

type Retweet struct {
}

func (c Retweet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	now := time.Now()

	// store that this user has liked this tweet

	// TODO(remy): user id
	if err := twitter.TweetDoneActionDAO().Retweet(user.Uid, ptid, now); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// rt on Twitter
	if _, err := twitter.GetAuthApi(user).Retweet(tid, true); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// plan an action to automatically unrt
	if au == "true" {
		unrt := twitter.UnRetweet{
			Uid:     uuid.New(),
			TweetId: ptid,
		}
		unrt.UserUid = user.Uid
		unrt.CreationTime = now
		unrt.ExecutionTime = now.Add(time.Hour * 24) // 1 day

		if err := unrt.Store(); err != nil {
			api.RenderErrJson(w, err)
			return
		}
	}

	api.RenderOk(w)
}

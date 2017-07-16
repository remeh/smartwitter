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

type Like struct {
}

func (c Like) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	now := time.Now()

	// store that this user has liked this tweet

	// TODO(remy): user id
	if err := twitter.TweetDoneActionDAO().Like(ptid, now); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// like on Twitter

	if _, err := twitter.GetAuthApi(user).Favorite(tid); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// plan an action to automatically unlike

	if au == "true" {
		unlike := twitter.UnLike{
			Uid:     uuid.New(),
			TweetId: ptid,
		}
		unlike.CreationTime = now
		unlike.ExecutionTime = now.Add(time.Hour * 144) // 6 days

		if err := unlike.Store(); err != nil {
			api.RenderErrJson(w, err)
			return
		}
	}

	api.RenderOk(w)
}

package action

import (
	"net/http"
	"time"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/twitter"
)

type Ignore struct {
}

func (c Ignore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if len(ptid) == 0 {
		api.RenderBadParameter(w, "tid")
		return
	}

	// ----------------------

	now := time.Now()

	// store that this user has ignored this tweet

	if err := twitter.TweetDoneActionDAO().Ignore(user.Uid, ptid, now); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderOk(w)
}

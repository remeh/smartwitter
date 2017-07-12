package action

import (
	"net/http"
	"time"

	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/twitter"
)

type Ignore struct {
}

func (c Ignore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// TODO(remy): user id
	if err := twitter.TweetDoneActionDAO().Ignore(ptid, now); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderOk(w)
}

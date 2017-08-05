package suggest

import (
	"net/http"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/suggest"
)

type Keywords struct{}

func (c Keywords) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// get keywords
	// ----------------------

	kwords, err := suggest.GetKeywords(user.Uid)
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderJson(w, 200, kwords)
}

package account

import (
	"net/http"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/api"
)

type Session struct{}

type sessionInfo struct {
}

func (c Session) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	api.RenderJson(w, 200, user)
}

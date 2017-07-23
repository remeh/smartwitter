package suggest

import (
	"net/http"
	"strconv"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/suggest"
)

type ChangeKeywords struct{}

func (c ChangeKeywords) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	// parse form
	// ----------------------

	r.ParseForm()
	keywords := api.Escape(r.Form["k"])
	pposition := api.EscapeOne(r.Form.Get("p"))
	position, err := strconv.Atoi(pposition)
	if err != nil || position < 0 {
		api.RenderBadParameter(w, "p")
		return
	}

	// set the keywords
	// ----------------------

	if err := suggest.SetKeywords(user.Uid, keywords, position); err != nil {
		api.RenderErrJson(w, err)
		return
	}
}

package adapter

import (
	"net/http"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/api"
)

type AuthHandler struct {
	handler http.Handler
}

func (a AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sc := api.GetSessionCookie(r)

	if len(sc) != 0 {
		if ok, err := account.UserDAO().Exists(sc); err != nil {
			api.RenderErrJson(w, err)
			return
		} else if !ok {
			api.RenderForbiddenJson(w)
			return
		}

		// ok
		sWriter := &StatusWriter{w, 200}
		a.handler.ServeHTTP(sWriter, r)
	}

	api.RenderForbiddenJson(w)
	return
}

// LogRoute creates a route which will log the route access.
func AuthAdapter(handler http.Handler) http.Handler {
	return AuthHandler{
		handler: handler,
	}
}

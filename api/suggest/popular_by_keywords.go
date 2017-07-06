package suggest

import (
	"net/http"

	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/suggest"
)

type ByKeywords struct {
}

func (c ByKeywords) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// parse form
	// ----------------------

	r.ParseForm()
	keywords := r.Form["k"]

	// TODO(remy): trim all keywords

	// get the suggestion
	// ----------------------

	tweets, err := suggest.SuggestByKeywords(keywords)

	// render the response
	// ----------------------

	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderJson(w, 200, tweets)
}

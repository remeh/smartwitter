package example

import (
	"net/http"

	"github.com/remeh/smartwitter/api"
)

type Example struct {
}

func (c Example) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.RenderOk(w)
}

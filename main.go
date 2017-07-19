package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/remeh/smartwitter/agent"
	"github.com/remeh/smartwitter/api/action"
	"github.com/remeh/smartwitter/api/adapter"
	"github.com/remeh/smartwitter/api/example"
	"github.com/remeh/smartwitter/api/suggest"
	"github.com/remeh/smartwitter/api/twitter"
	"github.com/remeh/smartwitter/config"
	l "github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
)

func main() {
	if _, err := storage.Init(config.Env().Conn); err != nil {
		l.Error("while init storage", err)
		os.Exit(1)
	}

	s := NewServer()
	declareApiRoutes(s)

	l.Info("Started.")

	// start the agents
	// ----------------------

	ctx, cf := context.WithTimeout(context.Background(), time.Minute*60)
	defer cf()

	go agent.GetTweets(ctx)
	go agent.PlannedActions(ctx)

	// start the webserver
	// ----------------------

	if err := s.Start(); err != nil {
		l.Error(err)
	}
}

func log(h http.Handler) http.Handler {
	return adapter.LogAdapter(h)
}

func auth(h http.Handler) http.Handler {
	return adapter.AuthAdapter(h)
}

func declareApiRoutes(s *Server) {
	s.AddApi("/example", log(example.Example{}), "GET")

	// ----------------------

	s.AddApi("/1.0/suggest", auth(log(suggest.ByKeywords{})), "GET")
	s.AddApi("/1.0/like", auth(log(action.Like{})), "POST")
	s.AddApi("/1.0/retweet", auth(log(action.Retweet{})), "POST")
	s.AddApi("/1.0/ignore", auth(log(action.Ignore{})), "POST")

	// twitter sign in
	// ----------------------

	if !config.Env().Debug {
		s.AddApi("/twitter/signin", log(twitter.RedirectUserToTwitter{}), "GET")
	} else {
		s.AddApi("/twitter/signin", log(twitter.DebugSignIn{}), "GET")
	}
	s.AddApi("/twitter/token", log(twitter.GetTwitterToken{}), "GET", "POST", "PUT")

	// debug route
	// ----------------------

}

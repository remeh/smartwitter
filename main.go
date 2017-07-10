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

	ctx, cf := context.WithTimeout(context.Background(), time.Second*15)
	defer cf()
	go agent.GetTweets(ctx)

	// start the webserver
	// ----------------------

	if err := s.Start(); err != nil {
		l.Error(err)
	}
}

func log(h http.Handler) http.Handler {
	return adapter.LogAdapter(h)
}

func startJobs() {
}

func declareApiRoutes(s *Server) {
	s.AddApi("/example", log(example.Example{}), "GET")

	// ----------------------

	s.AddApi("/1.0/suggest", log(suggest.ByKeywords{}), "GET")

	s.AddApi("/1.0/like", log(action.Like{}), "POST")
}

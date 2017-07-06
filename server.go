// Memoiz backend.
//
// Listening server.
//
// Rémy Mathieu © 2016

package main

import (
	"net/http"

	"github.com/remeh/smartwitter/config"
	l "github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	s := &Server{}
	s.router = mux.NewRouter()
	return s
}

// Starts listening.
func (s *Server) Start() error {
	conf := config.Env()

	// Opens the database connection.
	l.Info("opening the database connection.")
	_, err := storage.Init(conf.Conn)
	if err != nil {
		return err
	}

	// Prepares the router serving the static pages and assets.
	s.prepareStaticRouter()

	// Handles static routes
	http.Handle("/", s.router)

	// Starts listening.
	err = http.ListenAndServe(conf.ListenAddr, nil)
	return err
}

// AddApi adds a route in the API router of the application.
func (s *Server) AddApi(pattern string, handler http.Handler, methods ...string) {
	s.router.PathPrefix("/api").Subrouter().Handle(pattern, handler).Methods(methods...)
}

// ----------------------

func (s *Server) prepareStaticRouter() {
	conf := config.Env()
	// Add the final route, the static assets and pages.
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir(conf.PublicDir)))
	l.Info("serving static from directory", conf.PublicDir)
}

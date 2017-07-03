package main

import (
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/remeh/smartwitter/twitter"
)

func main() {
	v := url.Values{"track": os.Args[1:]}
	stream := twitter.GetApi().PublicStreamFilter(v)

	for s := range stream.C {
		t := s.(anaconda.Tweet)
		log.Println(t.User.Name, ": ", t.Text)
	}
}

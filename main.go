package main

import (
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func main() {
	c := EnvConfig()

	anaconda.SetConsumerKey(c.ConsumerKey)
	anaconda.SetConsumerSecret(c.ConsumerSecret)
	api := anaconda.NewTwitterApi(c.AccessToken, c.AccessTokenSecret)

	v := url.Values{"track": os.Args[1:]}
	stream := api.PublicStreamFilter(v)

	for s := range stream.C {
		t := s.(anaconda.Tweet)
		log.Println(t.User.Name, ": ", t.Text)
	}
}

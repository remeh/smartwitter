package main

import (
	"context"
	"os"
	"time"

	"github.com/remeh/smartwitter/agent"
	"github.com/remeh/smartwitter/config"
	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
)

func main() {
	if _, err := storage.Init(config.Env().Conn); err != nil {
		log.Error("while init storage", err)
		os.Exit(1)
	}

	log.Info("Started.")

	ctx, cf := context.WithTimeout(context.Background(), time.Second*15)
	defer cf()
	agent.GetTweets(ctx)
}

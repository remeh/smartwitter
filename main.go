package main

import (
	"context"
	"time"

	"github.com/remeh/smartwitter/agent"
	"github.com/remeh/smartwitter/config"
	"github.com/remeh/smartwitter/storage"
)

func main() {
	storage.Init(config.Env().Conn)

	ctx, cf := context.WithTimeout(context.Background(), time.Minute*5)
	defer cf()
	agent.GetTweets(ctx)
}

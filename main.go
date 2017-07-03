package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/remeh/smartwitter/follow"
	"github.com/remeh/smartwitter/twitter"
)

func main() {
	u, err := twitter.GetApi().GetUsersShow(os.Args[1], nil)
	if err != nil {
		log.Fatalln("err:", err)
	}
	id := u.Id
	println(id)

	ctx, cf := context.WithTimeout(context.Background(), time.Minute*600)
	defer cf()
	follow.Follow(ctx)
}

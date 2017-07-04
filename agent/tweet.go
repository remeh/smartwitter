package agent

import (
	"context"
	"log"
	"time"

	"github.com/remeh/smartwitter/twitter"
)

// GetTweets launches a crawler session.
// TODO(remy): better doc
func GetTweets(ctx context.Context) {
	for {
		after := time.After(time.Second * 3)

		// ----------------------

		select {
		case <-after:
			log.Println("debug: GetTweets is starting.")
			if err := run(ctx); err != nil {
				log.Println("error: while running GetTweets:", err)
			}
			log.Println("debug: GetTweets is ending.")
		case <-ctx.Done():
			log.Println("debug: GetTweets canceled.")
			return
		}
	}
}

func getTweets(ctx context.Context) error {
	sr, err := twitter.GetApi().GetSearch("golang", nil)
	if err != nil {
		return err
	}

	for _, s := range sr.Statuses {

		// store this tweet

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return nil
}

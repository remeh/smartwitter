package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/smartwitter/twitter"
)

// Unfavorite agent unfavoriting tweets when it's
// time to do so.
func Unfavorite(ctx context.Context) {
	for {
		after := time.After(time.Second * 3)

		// ----------------------

		select {
		case <-after:
			log.Debug("Unfavorite is starting.")
			if err := unfavorite(ctx); err != nil {
				log.Error("while running Unfavorite:", err)
			}
			log.Debug("Unfavorite is ending.")
		case <-ctx.Done():
			log.Debug("Unfavorite canceled.")
			return
		}
	}
}

// ----------------------

func unfavorite(ctx context.Context) error {
	var unfavs twitter.UnFavorites

	// find in database some tweets for which it's time to unfav.
	rows, err := storage.DB().Query(`
		SELECT "uid", "tweet_id"
		FROM "twitter_planned_action"
		WHERE
			"done" IS NULL
			and
			"execution_time" < now()
		LIMIT 20
	`)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var unfav twitter.UnFavorite

		if err := rows.Scan(&unfav.Uid, &unfav.TweetId); err != nil {
			return err
		}

		unfavs = append(unfavs, unfav)
	}

	// do not send all commands too fast
	errcount := 0
	for _, unfav := range unfavs {
		if _, err := twitter.GetApi().Unfavorite(unfav.TweetId); err != nil {
			log.Error(err)
			errcount++
			if errcount >= 3 {
				return fmt.Errorf("too much error unfavoriting ids. e.g. %v", err)
			}
			continue
		}

		if err := unfav.Forget(); err != nil {
			return err
		}

		log.Debug("unfaved", unfav.Uid)

		select {
		case <-ctx.Done():
			log.Debug("unfavorite canceled.")
			return ctx.Err()
		default:
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

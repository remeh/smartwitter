package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/smartwitter/twitter"
	"github.com/remeh/uuid"
)

// PlannedActions agent unfavoriting/unretweeting/...
// tweets when it's time to do so.
func PlannedActions(ctx context.Context) {
	for {
		after := time.After(time.Minute * 1)

		// ----------------------

		select {
		case <-after:
			log.Debug("PlannedActions is starting.")
			if err := plannedActions(ctx); err != nil {
				log.Error("while running PlannedActions:", err)
			}
			log.Debug("PlannedActions is ending.")
		case <-ctx.Done():
			log.Debug("PlannedActions canceled.")
			return
		}
	}
}

// ----------------------

func plannedActions(ctx context.Context) error {
	var actions twitter.PlannedActions

	// find in database some tweets for which it's time to unfav.
	rows, err := storage.DB().Query(`
		SELECT "type", "uid", "tweet_id"
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
		var typ, tid string
		var uid uuid.UUID

		if err := rows.Scan(&typ, &uid, &tid); err != nil {
			return err
		}

		var action twitter.PlannedAction

		switch typ {
		case "unretweet":
			action = &twitter.UnRetweet{
				Uid:     uid,
				TweetId: tid,
			}
		case "unlike":
			action = &twitter.UnLike{
				Uid:     uid,
				TweetId: tid,
			}
		}

		actions = append(actions, action)
	}

	// do not send all commands too fast
	errcount := 0
	for _, action := range actions {
		if err := action.Do(); err != nil {
			log.Error(err)
			errcount++
			if errcount >= 3 {
				return fmt.Errorf("too much error in planned actions agent. e.g. %v", err)
			}
			continue
		}

		if err := action.Forget(); err != nil {
			return err
		}

		log.Debug(action.String())

		select {
		case <-ctx.Done():
			log.Debug("plannedActions canceled.")
			return ctx.Err()
		default:
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

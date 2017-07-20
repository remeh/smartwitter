package agent

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/twitter"
	"github.com/remeh/uuid"
)

// GetTweets launches a crawler session.
// TODO(remy): better doc
func GetTweets(ctx context.Context) {
	for {
		after := time.After(time.Minute * 5)

		// ----------------------

		select {
		case <-after:
			log.Debug("GetTweets is starting.")

			k, err := getNextKeywords()
			if err != nil {
				log.Error("GetTweets: getNextKeywords:", err)
				continue
			}

			if err = getTweets(ctx, k); err != nil {
				log.Error("while running GetTweets:", err)
			}

			log.Debug("GetTweets is ending.")
		case <-ctx.Done():
			log.Debug("GetTweets canceled.")
			return
		}
	}
}

// ----------------------

type keywordsToSearch struct {
	Uid      uuid.UUID
	UserUid  uuid.UUID
	Keywords string
}

func getNextKeywords() (keywordsToSearch, error) {
	// TODO(remy):
	return keywordsToSearch{}, nil
}

func getTweets(ctx context.Context, k keywordsToSearch) error {
	v := url.Values{
		"tweet_mode":  []string{"extended"},
		"lang":        []string{"en"},
		"count":       []string{"50"},
		"result_type": []string{"mixed"},
	}

	sr, err := twitter.GetApi().GetSearch("golang code -filter:retweets", v)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, s := range sr.Statuses {
		// atm, we want to ignore the retweets to
		// not fill up the database with retweets.
		if strings.Contains(s.FullText, "RT @") {
			continue
		}

		// create this tweet and twitter user
		// ----------------------

		// tweet

		t := twitter.TweetFromTweet(s, now, []string{"golang", "code"})

		// twitter user

		tu := twitter.TwitterUserFromTweet(s, now)

		// upsert
		// ----------------------

		if err := twitter.TwitterUserDAO().Upsert(tu); err != nil {
			return log.Err("getTweets: upsert TwitterUser:", err)
		}

		if err := twitter.TweetDAO().Upsert(t); err != nil {
			return log.Err("getTweets: upsert Tweet:", err)
			return err
		}

		log.Debug("stored tweet:", tu.Name, t.Text)

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return nil
}

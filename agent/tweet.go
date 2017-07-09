package agent

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/remeh/smartwitter/log"
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
			log.Debug("GetTweets is starting.")
			if err := getTweets(ctx); err != nil {
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

func getTweets(ctx context.Context) error {
	v := url.Values{
		"tweet_mode":  []string{"extended"},
		"lang":        []string{"en"},
		"count":       []string{"30"},
		"result_type": []string{"recent"},
	}
	sr, err := twitter.GetApi().GetSearch("golang", v)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, s := range sr.Statuses {
		// atm, we want to ignore the retweets to
		// not fill up the database with retweets.
		if strings.Contains(s.FullText, "RT @") {
			log.Debug("ignoring tweet:", s.User.Name, s.FullText)
			continue
		}

		// create this tweet and twitter user
		// ----------------------

		// tweet

		t := twitter.TweetFromTweet(s, now, []string{"golang"})

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

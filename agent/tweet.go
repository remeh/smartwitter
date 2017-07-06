package agent

import (
	"context"
	"fmt"
	"net/url"
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
		"count":       []string{"50"},
		"result_type": []string{"recent"},
	}
	sr, err := twitter.GetApi().GetSearch("golang", v)
	if err != nil {
		return err
	}

	for _, s := range sr.Statuses {
		now := time.Now()

		// create this tweet and twitter user
		// ----------------------

		// tweet

		// TODO(remy): export this method into package twiter

		t := twitter.NewTweet(s.Id, s.User.Id)
		if t.TwitterCreationTime, err = s.CreatedAtTime(); err != nil {
			log.Warning("getTweets: getting tweet creation time:", err)
		}
		t.CreationTime = now
		t.LastUpdate = now
		t.Text = s.FullText
		t.RetweetCount = s.RetweetCount
		t.FavoriteCount = s.FavoriteCount
		t.Lang = s.Lang
		t.Keywords = []string{"golang"}
		t.Link = fmt.Sprintf("https://twitter.com/%s/status/%d", s.User.ScreenName, s.Id)

		// twitter user

		tu := twitter.NewTwitterUser(s.User.Id)
		tu.CreationTime = now
		tu.LastUpdate = now
		tu.Description = s.User.Description
		tu.Name = s.User.Name
		tu.ScreenName = s.User.ScreenName
		tu.TimeZone = s.User.TimeZone
		tu.UtcOffset = s.User.UtcOffset
		tu.FollowersCount = s.User.FollowersCount

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

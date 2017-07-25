package agent

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/smartwitter/twitter"
	"github.com/remeh/uuid"

	"github.com/lib/pq"
)

var (
	// Do not crawl too many times the keywords.
	// For each keywords, this time should be elapsed
	// to be crawled again.
	MinIntervalBetweenGetTweets = "5 minutes"
)

// GetTweets launches a crawler session.
// TODO(remy): better doc
func GetTweets(ctx context.Context) {
	for {
		after := time.After(time.Second * 1)

		// ----------------------

		select {
		case <-after:
			// TODO(remy: ctx ?

			k, err := getNextKeywords()
			if err != nil {
				log.Error(err)
				break
			}

			if k == nil {
				break
			}

			if err = getTweets(ctx, k); err != nil {
				log.Error(err)
			}

			if err = updateTime(k); err != nil {
				log.Error(err)
			}

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
	Keywords []string
}

func getNextKeywords() (*keywordsToSearch, error) {
	kts := keywordsToSearch{}

	if err := storage.DB().QueryRow(`
		select
			"uid", "user_uid", "keywords"
		from "twitter_keywords_watcher"
		where
			coalesce("last_run", '1970-01-01')
				+ interval '`+MinIntervalBetweenGetTweets+`' < now()
		order by coalesce("last_run", '1970-01-01')
		limit 1
	`).Scan(
		&kts.Uid,
		&kts.UserUid,
		pq.Array(&kts.Keywords),
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, log.Err("getNextKeywords", err)
		}
	}

	return &kts, nil
}

func updateTime(k *keywordsToSearch) error {
	if k == nil {
		return fmt.Errorf("updateTime: called with nil object")
	}

	if _, err := storage.DB().Exec(`
		UPDATE "twitter_keywords_watcher"
		SET
			"last_run" = now()
		WHERE
			"uid" = $1
	`, k.Uid); err != nil {
		return log.Err("updateTime", err)
	}
	return nil
}

func getTweets(ctx context.Context, k *keywordsToSearch) error {
	if k == nil {
		return fmt.Errorf("getTweets: called with nil object")
	}

	v := url.Values{
		"tweet_mode":  []string{"extended"},
		"lang":        []string{"en"}, // TODO(remy): user language
		"count":       []string{"25"},
		"result_type": []string{"mixed"},
	}

	if k.Uid.IsNil() || k.UserUid.IsNil() {
		return fmt.Errorf("getTweets called with nil k.Uid or k.UserUid")
	}

	// Get the user
	user, err := account.UserDAO().Find(k.UserUid)
	if err != nil {
		return log.Err("getTweets: while looking for the user", err)
	}

	str := fmt.Sprintf("%s -filter:retweets", strings.Join(k.Keywords, " "))
	sr, err := twitter.GetAuthApi(user).GetSearch(str, v)
	if err != nil {
		return log.Err("getTweets: while calling api", err)
	}

	stored := 0
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

		t := twitter.TweetFromTweet(s, now, k.Keywords)

		// twitter user

		tu := twitter.TwitterUserFromTweet(s, now)

		// upsert
		// ----------------------

		if err := twitter.TwitterUserDAO().Upsert(tu); err != nil {
			return log.Err("getTweets: upsert TwitterUser:", err)
		}

		if err := twitter.TweetDAO().Upsert(t); err != nil {
			return log.Err("getTweets: upsert Tweet:", err)
		}

		log.Debug("stored tweet:", tu.Name, t.Text)
		stored++

		select {
		case <-ctx.Done():
			return log.Err("getTweets", ctx.Err())
		default:
		}
	}

	log.Info("For keywords", k.Keywords, "upserted", stored, "tweets")

	return nil
}

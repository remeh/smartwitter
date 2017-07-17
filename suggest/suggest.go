package suggest

import (
	"database/sql"
	"time"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/smartwitter/twitter"

	"github.com/lib/pq"
)

// TODO(remy): user
func SuggestByKeywords(user *account.User, keywords []string, duration time.Duration, limit int) (twitter.Tweets, twitter.TweetDoneActions, error) {
	var rows *sql.Rows
	var tids []string
	var err error

	// get tweets with these keywords
	// ----------------------

	ct := time.Now().Add(-duration)

	if rows, err = storage.DB().Query(`
		select text, tweet.twitter_id, retweet_count, favorite_count, link, twitter_user_uid
		from "tweet"
		join "twitter_user" on
			"tweet"."twitter_user_uid" = "twitter_user"."uid"
		where
			not "text" LIKE 'RT @%'
			and
			"tweet"."keywords" @> $1::text[]
			and
			"tweet"."creation_time" > $2
		order by (favorite_count+retweet_count+followers_count) desc, "tweet".uid
		limit $3;
	`, pq.Array(keywords), ct, limit); err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	rv := make(twitter.Tweets, 0)
	for rows.Next() {
		t := twitter.Tweet{}

		if err := rows.Scan(&t.Text, &t.TwitterId, &t.RetweetCount, &t.FavoriteCount, &t.Link, &t.TwitterUserUid); err != nil {
			return nil, nil, err
		}

		tids = append(tids, t.TwitterId)

		rv = append(rv, t)
	}

	// tweet states for this user
	// ----------------------

	tdas, err := twitter.TweetDoneActionDAO().FindByTweets(user.Uid, tids)
	if err != nil {
		return nil, nil, err
	}

	return rv, tdas, nil
}

package suggest

import (
	"database/sql"
	"time"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/smartwitter/twitter"

	"github.com/lib/pq"
)

func SuggestByKeywords(user *account.User, keywords []string, duration time.Duration, limit int) (twitter.Tweets, twitter.TweetDoneActions, error) {
	var rows *sql.Rows
	var err error

	// get tweets with these keywords
	// ----------------------

	ct := time.Now().Add(-duration)

	if rows, err = storage.DB().Query(`
		select
			"tweet"."text",
			"tweet"."twitter_id",
			"tweet"."retweet_count",
			"tweet"."favorite_count",
			"tweet"."link",
			"tweet"."twitter_user_uid",
			"tweet_done_action"."ignored_time",
			"tweet_done_action"."liked_time",
			"tweet_done_action"."retweeted_time",
			"tweet_done_action"."ignored_time" IS NOT NULL,
			"tweet_done_action"."liked_time" IS NOT NULL,
			"tweet_done_action"."retweeted_time" IS NOT NULL
		from "tweet"
		join "twitter_user" on
			"tweet"."twitter_user_uid" = "twitter_user"."uid"
		left join "tweet_done_action" on
			"tweet_done_action"."user_uid" = $3
			and
			"tweet"."twitter_id" = "tweet_done_action"."tweet_id"
		where
			not "text" LIKE 'RT @%'
			and
			"tweet"."keywords" @> $1::text[]
			and
			"tweet"."creation_time" > $2
			and
			"tweet_done_action"."ignored_time" IS NULL
		order by (favorite_count+retweet_count) desc, "tweet".uid
		limit $4;
	`, pq.Array(keywords), ct, user.Uid, limit); err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	rv := make(twitter.Tweets, 0)
	tdas := make(twitter.TweetDoneActions, 0)

	for rows.Next() {
		t := twitter.Tweet{}
		tda := twitter.TweetDoneAction{}

		var it, lt, rt pq.NullTime

		if err := rows.Scan(
			&t.Text,
			&t.TwitterId,
			&t.RetweetCount,
			&t.FavoriteCount,
			&t.Link,
			&t.TwitterUserUid,
			&it,
			&lt,
			&rt,
			&tda.Ignored,
			&tda.Liked,
			&tda.Retweeted,
		); err != nil {
			return nil, nil, err
		}

		if it.Valid {
			tda.IgnoreTime = it.Time
		}
		if lt.Valid {
			tda.LikeTime = lt.Time
		}
		if rt.Valid {
			tda.RetweetTime = rt.Time
		}

		rv = append(rv, t)
		tdas = append(tdas, tda)
	}

	return rv, tdas, nil
}

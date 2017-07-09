package suggest

import (
	"database/sql"

	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/smartwitter/twitter"

	"github.com/lib/pq"
)

func SuggestByKeywords(keywords []string) (twitter.Tweets, error) {
	var rows *sql.Rows
	var err error

	if rows, err = storage.DB().Query(`
		select text, tweet.twitter_id, retweet_count, favorite_count, link, twitter_user_uid
		from "tweet"
		join "twitter_user" on
			"tweet"."twitter_user_uid" = "twitter_user"."uid"
		where
			not "text" LIKE 'RT @%'
			and
			"tweet"."keywords" <@ $1::text[]
		order by (favorite_count+retweet_count+followers_count) desc, "tweet".uid
		limit 3;
	`, pq.Array(keywords)); err != nil {
		return nil, err
	}

	defer rows.Close()

	rv := make(twitter.Tweets, 0)
	for rows.Next() {
		t := twitter.Tweet{}

		if err := rows.Scan(&t.Text, &t.TwitterId, &t.RetweetCount, &t.FavoriteCount, &t.Link, &t.TwitterUserUid); err != nil {
			return nil, err
		}

		rv = append(rv, t)
	}

	return rv, nil
}

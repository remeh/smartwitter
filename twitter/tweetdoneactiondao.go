package twitter

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
)

type tweetDoneActionDAO struct {
	DB *sql.DB
}

// ----------------------

var tdaDAO *tweetDoneActionDAO

func TweetDoneActionDAO() *tweetDoneActionDAO {
	if dao != nil {
		return tdaDAO
	}

	tdaDAO = &tweetDoneActionDAO{
		DB: storage.DB(),
	}

	if err := tdaDAO.InitStmt(); err != nil {
		log.Error("Can't prepare TweetDoneActionDAO")
		panic(err)
	}

	return tdaDAO
}

func (d *tweetDoneActionDAO) InitStmt() error {
	var err error
	return err
}

// TODO(remy): user
func (d *tweetDoneActionDAO) Like(tid string, t time.Time) error {
	if _, err := storage.DB().Exec(`
		INSERT INTO "tweet_done_action"
		("user_uid", "tweet_id", "liked_time")
		VALUES
		($1, $2, $3)
		ON CONFLICT ("user_uid", "tweet_id") DO UPDATE SET
			"liked_time" = $3
	`, "", tid, t); err != nil {
		return err
	}
	return nil
}

// TODO(remy): user
func (d *tweetDoneActionDAO) Retweet(tid string, t time.Time) error {
	if _, err := storage.DB().Exec(`
		INSERT INTO "tweet_done_action"
		("user_uid", "tweet_id", "retweeted_time")
		VALUES
		($1, $2, $3)
		ON CONFLICT ("user_uid", "tweet_id") DO UPDATE SET
			"retweeted_time" = $3
	`, "", tid, t); err != nil {
		return err
	}
	return nil
}

// TODO(remy): user
func (d *tweetDoneActionDAO) Ignore(tid string, t time.Time) error {
	if _, err := storage.DB().Exec(`
		INSERT INTO "tweet_done_action"
		("user_uid", "tweet_id", "ignore_time")
		VALUES
		($1, $2, $3)
		ON CONFLICT ("user_uid", "tweet_id") DO UPDATE SET
			"ignore_time" = $3
	`, "", tid, t); err != nil {
		return err
	}
	return nil
}

func (d *tweetDoneActionDAO) FindByTweets(tids []string) (TweetDoneActions, error) {
	params := make([]interface{}, len(tids))
	for i := 0; i < len(tids); i++ {
		params[i] = tids[i]
	}

	rows, err := storage.DB().Query(`
		select distinct
			"tweet"."twitter_id",
			"tweet_done_action"."ignored_time",
			"tweet_done_action"."liked_time",
			"tweet_done_action"."retweeted_time",
			"tweet_done_action"."ignored_time" IS NOT NULL,
			"tweet_done_action"."liked_time" IS NOT NULL,
			"tweet_done_action"."retweeted_time" IS NOT NULL
		from "tweet"
		left join "tweet_done_action" ON
			"tweet_done_action"."tweet_id" = "tweet"."twitter_id"
			-- TODO and user =
		where
			"tweet"."twitter_id" IN `+storage.InClause(1, len(tids))+`
		`, params...)

	if err != nil {
		return nil, err
	}

	defer rows.Next()

	var tdas TweetDoneActions

	for rows.Next() {
		var tda TweetDoneAction

		var ignoreTime, likeTime, rtTime pq.NullTime

		if err := rows.Scan(&tda.TweetId, &ignoreTime, &likeTime, &rtTime, &tda.Ignored, &tda.Liked, &tda.Retweeted); err != nil {
			return nil, err
		}

		if ignoreTime.Valid {
			tda.IgnoreTime = ignoreTime.Time
		}

		if likeTime.Valid {
			tda.LikeTime = likeTime.Time
		}

		if rtTime.Valid {
			tda.RetweetTime = rtTime.Time
		}

		tdas = append(tdas, tda)
	}

	return tdas, nil
}

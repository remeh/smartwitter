package twitter

import (
	"database/sql"

	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"

	"github.com/lib/pq"
)

type tweetDAO struct {
	DB *sql.DB
}

// ----------------------

var dao *tweetDAO

func TweetDAO() *tweetDAO {
	if dao != nil {
		return dao
	}

	dao = &tweetDAO{
		DB: storage.DB(),
	}

	if err := dao.InitStmt(); err != nil {
		log.Error("Can't prepare TweetDAO")
		panic(err)
	}

	return dao
}

func (d *tweetDAO) InitStmt() error {
	var err error
	return err
}

func (d *tweetDAO) Upsert(t *Tweet) error {
	if _, err := d.DB.Exec(`
		INSERT INTO "tweet" ("uid", "creation_time", "last_update", "twitter_id", "twitter_creation_time", "text", "twitter_user_uid", "retweet_count", "favorite_count", "lang", "keywords", "link")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT ("uid") DO UPDATE SET
			"last_update" = $3,
			"twitter_id" = $4,
			"twitter_creation_time" = $5,
			"text" = $6,
			"twitter_user_uid" = $7,
			"retweet_count" = $8,
			"favorite_count" = $9,
			"lang" = $10,
			"keywords" = array(select distinct unnest("tweet"."keywords" || $11)),
			"link" = $12
	`, t.Uid(), t.CreationTime, t.LastUpdate, t.TwitterId, t.TwitterCreationTime, t.Text, t.TwitterUserUid, t.RetweetCount, t.FavoriteCount, t.Lang, pq.Array(t.Keywords), t.Link); err != nil {
		return err
	}
	return nil
}

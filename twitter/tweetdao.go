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
	idxStart, idxEnd := t.Entities.Indices()
	if _, err := d.DB.Exec(`
		INSERT INTO "tweet" ("uid", "creation_time", "last_update", "twitter_id", "twitter_creation_time", "text", "twitter_user_uid", "retweet_count", "favorite_count", "lang", "keywords", "link", "entities_type", "entities_display_url", "entities_url", "entities_idx_start", "entities_idx_end", "entities_screen_name", "entities_user_name", "entities_user_id", "entities_hashtag")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
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
			"link" = $12,
			"entities_type" = $13,
			"entities_display_url" = $14,
			"entities_url" = $15,
			"entities_idx_start" = $16,
			"entities_idx_end" = $17,
			"entities_screen_name" = $18,
			"entities_user_name" = $19,
			"entities_user_id" = $20,
			"entities_hashtag" = $21
	`, t.Uid(), t.CreationTime, t.LastUpdate, t.TwitterId, t.TwitterCreationTime, t.Text, t.TwitterUserUid, t.RetweetCount, t.FavoriteCount, t.Lang, pq.Array(t.Keywords), t.Link,
		pq.Array(t.Entities.Types()),
		pq.Array(t.Entities.DisplayUrls()),
		pq.Array(t.Entities.Urls()),
		pq.Array(idxStart),
		pq.Array(idxEnd),
		pq.Array(t.Entities.ScreenNames()),
		pq.Array(t.Entities.UserNames()),
		pq.Array(t.Entities.UserTwitterIds()),
		pq.Array(t.Entities.Hashtags()),
	); err != nil {
		return err
	}
	return nil
}

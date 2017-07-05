package twitter

import (
	"database/sql"

	"remy.io/memoiz/log"
	"remy.io/memoiz/storage"
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

func (d *tweetDAO) UpsertTweet(t *Tweet) error {
	if _, err := d.DB.Exec(`
		INSERT INTO "tweet" ("uid", "creation_time", "twitter_id", "twitter_creation_time", "text", "user_uid")
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT ("uid") DO UPDATE SET
			"creation_time" = $2,
			"twitter_id" = $3,
			"twitter_creation_time" = $4,
			"text" = $5,
			"user_uid" = $6
	`, t.Uid(), t.TwitterId, t.TwitterCreationTime, t.Text, t.UserUid); err != nil {
		return err
	}
	return nil
}

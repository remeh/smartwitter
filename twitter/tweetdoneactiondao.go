package twitter

import (
	"database/sql"
	"time"

	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/uuid"
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

func (d *tweetDoneActionDAO) Like(userUid uuid.UUID, tid string, t time.Time) error {
	if _, err := storage.DB().Exec(`
		INSERT INTO "tweet_done_action"
		("user_uid", "tweet_id", "liked_time")
		VALUES
		($1, $2, $3)
		ON CONFLICT ("user_uid", "tweet_id") DO UPDATE SET
			"liked_time" = $3
	`, userUid, tid, t); err != nil {
		return err
	}
	return nil
}

func (d *tweetDoneActionDAO) Retweet(userUid uuid.UUID, tid string, t time.Time) error {
	if _, err := storage.DB().Exec(`
		INSERT INTO "tweet_done_action"
		("user_uid", "tweet_id", "retweeted_time")
		VALUES
		($1, $2, $3)
		ON CONFLICT ("user_uid", "tweet_id") DO UPDATE SET
			"retweeted_time" = $3
	`, userUid, tid, t); err != nil {
		return err
	}
	return nil
}

func (d *tweetDoneActionDAO) Ignore(userUid uuid.UUID, tid string, t time.Time) error {
	if _, err := storage.DB().Exec(`
		INSERT INTO "tweet_done_action"
		("user_uid", "tweet_id", "ignored_time")
		VALUES
		($1, $2, $3)
		ON CONFLICT ("user_uid", "tweet_id") DO UPDATE SET
			"ignored_time" = $3
	`, userUid, tid, t); err != nil {
		return err
	}
	return nil
}

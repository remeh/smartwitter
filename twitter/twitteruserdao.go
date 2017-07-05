package twitter

import (
	"database/sql"
	"fmt"

	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/uuid"
)

type twitterUserDAO struct {
	DB *sql.DB
}

// ----------------------

var tuDao *twitterUserDAO

func TwitterUserDAO() *twitterUserDAO {
	if tuDao != nil {
		return tuDao
	}

	tuDao = &twitterUserDAO{
		DB: storage.DB(),
	}

	if err := tuDao.InitStmt(); err != nil {
		log.Error("Can't prepare TwitterUserDAO")
		panic(err)
	}

	return tuDao
}

func (d *twitterUserDAO) InitStmt() error {
	var err error
	return err
}

func (d *twitterUserDAO) Upsert(tu *TwitterUser) error {
	if tu == nil {
		return fmt.Errorf("tu == nil")
	}

	if _, err := d.DB.Exec(`
		INSERT INTO "tweet_user" ("uid", "creation_time", "twitter_id", "twitter_creation_time", "description", "screen_name", "name", "timezone", "utc_offset")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT ("uid") DO UPDATE SET
			"creation_time" = $2,
			"twitter_id" = $3,
			"description" = $4,
			"screen_name" = $5,
			"name" = $6,
			"timezone" = $7,
			"utc_offset" = $8
	`, tu.Uid(), tu.CreationTime, tu.TwitterId, tu.Description, tu.ScreenName, tu.Name, tu.TimeZone, tu.UtcOffset); err != nil {
		return err
	}

	return nil
}

func (d *twitterUserDAO) Find(id uuid.UUID) (*TwitterUser, error) {
	if id == nil {
		return nil, fmt.Errorf("nil db or nil id")
	}

	rv := &TwitterUser{uid: id}
	if err := d.DB.QueryRow(`
		SELECT "creation_time", "twitter_id", "description", "screen_name", "name", "timezone", "utc_offset" FROM "tweet_user"
		WHERE "uid" = $1
		LIMIT 1
	`).Scan(&rv.CreationTime, &rv.TwitterId, &rv.Description, &rv.ScreenName, &rv.Name, &rv.TimeZone, &rv.UtcOffset); err != nil {
		return nil, err
	}
	return rv, nil
}

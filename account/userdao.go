package account

import (
	"database/sql"

	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
)

type userDAO struct {
	DB *sql.DB
}

// ----------------------

var dao *userDAO

func UserDAO() *userDAO {
	if dao != nil {
		return dao
	}

	dao = &userDAO{
		DB: storage.DB(),
	}

	if err := dao.InitStmt(); err != nil {
		log.Error("Can't prepare UserDAO")
		panic(err)
	}

	return dao
}

func (d *userDAO) InitStmt() error {
	var err error
	return err
}

func (d *userDAO) UpsertOnLogin(u *User) error {
	if _, err := d.DB.Exec(`
		INSERT INTO "user" ("uid", "creation_time", "last_login", "twitter_token", "twitter_secret", "twitter_id", "twitter_name", "twitter_username")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT ("uid") DO UPDATE SET
			"creation_time" = $2,
			"last_login" = $3,
			"twitter_token" = $4,
			"twitter_secret" = $5,
			"twitter_id" = $6,
			"twitter_name" = $7,
			"twitter_username" = $8
	`, u.Uid, u.CreationTime, u.LastLogin, u.TwitterToken, u.TwitterSecret, u.TwitterId, u.TwitterName, u.TwitterUsername); err != nil {
		return err
	}
	return nil
}

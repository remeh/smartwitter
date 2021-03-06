package account

import (
	"database/sql"

	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/uuid"
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
		INSERT INTO "user" ("uid", "creation_time", "last_login", "twitter_token", "twitter_secret", "twitter_id", "twitter_name", "twitter_username", "session_token")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT ("uid") DO UPDATE SET
			"creation_time" = $2,
			"last_login" = $3,
			"twitter_token" = $4,
			"twitter_secret" = $5,
			"twitter_id" = $6,
			"twitter_name" = $7,
			"twitter_username" = $8,
			"twitter_avatar" = $9,
			"session_token" = $10
	`, u.Uid, u.CreationTime, u.LastLogin, u.TwitterToken, u.TwitterSecret, u.TwitterId, u.TwitterName, u.TwitterUsername, u.TwitterAvatar, u.SessionToken); err != nil {
		return err
	}
	return nil
}

func (d *userDAO) Find(id uuid.UUID) (*User, error) {
	rv := &User{}

	if err := d.DB.QueryRow(`
		SELECT "uid", "creation_time", "last_login", "twitter_token", "twitter_secret", "twitter_id", "twitter_name", "twitter_username", "twitter_avatar", "session_token" FROM "user"
		WHERE
			"uid" = $1
		LIMIT 1
	`, id).Scan(
		&rv.Uid,
		&rv.CreationTime,
		&rv.LastLogin,
		&rv.TwitterToken,
		&rv.TwitterSecret,
		&rv.TwitterId,
		&rv.TwitterName,
		&rv.TwitterUsername,
		&rv.TwitterAvatar,
		&rv.SessionToken); err != nil {
		return nil, err
	}

	return rv, nil
}

func (d *userDAO) FindBySession(sessionToken string) (*User, error) {
	rv := &User{}

	if err := d.DB.QueryRow(`
		SELECT "uid", "creation_time", "last_login", "twitter_token", "twitter_secret", "twitter_id", "twitter_name", "twitter_username", "twitter_avatar" FROM "user"
		WHERE
			session_token = $1
		LIMIT 1
	`, sessionToken).Scan(
		&rv.Uid,
		&rv.CreationTime,
		&rv.LastLogin,
		&rv.TwitterToken,
		&rv.TwitterSecret,
		&rv.TwitterId,
		&rv.TwitterName,
		&rv.TwitterUsername,
		&rv.TwitterAvatar); err != nil {
		return nil, err
	}

	rv.SessionToken = sessionToken

	return rv, nil
}

func (d *userDAO) Exists(sessionToken string) (bool, error) {
	var s int
	if err := d.DB.QueryRow(`
		SELECT length("session_token") FROM "user"
		WHERE
			"session_token" = $1
			AND
			"session_token" IS NOT NULL
		LIMIT 1
	`, sessionToken).Scan(&s); err != nil {
		return false, err
	}
	return s > 0, nil
}

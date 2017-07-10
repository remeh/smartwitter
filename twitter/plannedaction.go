package twitter

import (
	"fmt"
	"time"

	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/uuid"
)

type PlannedActions []PlannedAction

type PlannedAction interface {
	// Store the action in database
	Store() error
	// Forget the action in database
	Forget() error
	// Apply the action
	Do() error
}

// ----------------------

type action struct {
	// TODO(remy): add user
	// Time at which this action has been created
	CreationTime time.Time
	// Time at which the action must be executed
	ExecutionTime time.Time
}

type UnRetweets []UnRetweet

type UnRetweet struct {
	Uid     uuid.UUID
	TweetId int64
	action
}

func (u *UnRetweet) Do() error {
	_, err := GetApi().UnRetweet(u.TweetId, true)
	return err
}

type UnFavorites []UnFavorite

type UnFavorite struct {
	Uid     uuid.UUID
	TweetId int64
	action
}

func (u *UnFavorite) Do() error {
	_, err := GetApi().Unfavorite(u.TweetId)
	return err
}

func (u *UnFavorite) Store() error {
	if u.Uid.IsNil() {
		return fmt.Errorf("UnFavorite.Store(): nil Uid")
	}

	_, err := storage.DB().Exec(`
		INSERT INTO "twitter_planned_action"
		("uid", "type", "tweet_id", "creation_time", "execution_time", "done")
		VALUES
		($1, 'unfavorite', $2, $3, $4, NULL)
		`, u.Uid, u.TweetId, u.CreationTime, u.ExecutionTime)

	return err
}

func (u *UnFavorite) Forget() error {
	if u.Uid.IsNil() {
		return fmt.Errorf("UnFavorite.Forget(): nil Uid")
	}

	_, err := storage.DB().Exec(`
		UPDATE "twitter_planned_action"
		SET
			"done" = now()
		WHERE
			"uid" = $1
	`, u.Uid)
	return err
}

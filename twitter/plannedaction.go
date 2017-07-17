package twitter

import (
	"fmt"
	"strconv"
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
	String() string
}

type action struct {
	UserUid uuid.UUID

	// Time at which this action has been created
	CreationTime time.Time
	// Time at which the action must be executed
	ExecutionTime time.Time
}

// ----------------------

type UnRetweets []UnRetweet

type UnRetweet struct {
	Uid     uuid.UUID
	TweetId string
	action
}

func (u *UnRetweet) Do() error {
	tid, err := strconv.ParseInt(u.TweetId, 10, 64)
	if err != nil {
		return err
	}

	_, err = GetApi().UnRetweet(tid, true)
	return err
}

func (u *UnRetweet) Store() error {
	if u.Uid.IsNil() {
		return fmt.Errorf("UnRetweet.Store(): nil Uid")
	}

	_, err := storage.DB().Exec(`
		INSERT INTO "twitter_planned_action"
		("uid", "user_uid", "type", "tweet_id", "creation_time", "execution_time", "done")
		VALUES
		($1, $2, 'unretweet', $3, $4, $5, NULL)
		`, u.Uid, u.UserUid, u.TweetId, u.CreationTime, u.ExecutionTime)

	return err
}

func (u *UnRetweet) Forget() error {
	if u.Uid.IsNil() {
		return fmt.Errorf("UnRetweet.Forget(): nil Uid")
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

func (u *UnRetweet) String() string {
	return fmt.Sprintf("UnRetweet: %v", u.Uid)
}

// ----------------------

type UnLikes []UnLike

type UnLike struct {
	Uid     uuid.UUID
	TweetId string
	action
}

func (u *UnLike) Do() error {
	tid, err := strconv.ParseInt(u.TweetId, 10, 64)
	if err != nil {
		return err
	}

	_, err = GetApi().Unfavorite(tid)
	return err
}

func (u *UnLike) Store() error {
	if u.Uid.IsNil() {
		return fmt.Errorf("UnLike.Store(): nil Uid")
	}

	_, err := storage.DB().Exec(`
		INSERT INTO "twitter_planned_action"
		("uid", "user_uid", "type", "tweet_id", "creation_time", "execution_time", "done")
		VALUES
		($1, $2, 'unlike', $3, $4, $5, NULL)
		`, u.Uid, u.UserUid, u.TweetId, u.CreationTime, u.ExecutionTime)

	return err
}

func (u *UnLike) Forget() error {
	if u.Uid.IsNil() {
		return fmt.Errorf("UnLike.Forget(): nil Uid")
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
func (u *UnLike) String() string {
	return fmt.Sprintf("UnLike: %v", u.Uid)
}

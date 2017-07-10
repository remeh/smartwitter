package twitter

import "github.com/lib/pq"
import "github.com/lib/pq/hstore"
import "time"

type FutureActions []FutureAction

type FutureAction interface {
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

type UnRetweet struct {
	action
	TweetId int64
}

func (u *UnRetweet) Do() error {
	_, err := GetApi().UnRetweet(u.TweetId, true)
	return err
}

type UnFavorite struct {
	action
	TweetId int64
}

func (u *UnFavorite) Do() error {
	_, err := GetApi().Unfavorite(u.TweetId)
	return err
}

func (u *UnFavorite) Store() error {
	panic("todo")
}

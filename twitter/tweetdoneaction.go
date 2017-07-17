package twitter

import (
	"time"

	"github.com/remeh/uuid"
)

type TweetDoneActions []TweetDoneAction

type TweetDoneAction struct {
	TweetId     string
	UserUid     uuid.UUID
	Ignored     bool
	Retweeted   bool
	Liked       bool
	IgnoreTime  time.Time
	RetweetTime time.Time
	LikeTime    time.Time
}

func (tws TweetDoneActions) Get(tid string) (TweetDoneAction, bool) {
	for _, tw := range tws {
		if tw.TweetId == tid {
			return tw, true
		}
	}
	return TweetDoneAction{TweetId: tid}, false
}

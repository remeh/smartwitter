package twitter

import (
	"time"

	"github.com/remeh/smartwitter/log"
)

type TweetDoneActions []TweetDoneAction

type TweetDoneAction struct {
	TweetId     string
	Retweeted   bool
	Liked       bool
	RetweetTime time.Time
	LikeTime    time.Time
}

func (tws TweetDoneActions) Get(tid string) TweetDoneAction {
	for _, tw := range tws {
		if tw.TweetId == tid {
			return tw
		}
	}
	log.Warning("TweetDoneActions.Get() didn't find the TweetDoneAction.")
	return TweetDoneAction{TweetId: tid}
}

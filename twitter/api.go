package twitter

import (
	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/config"

	"github.com/remeh/anaconda"
)

var api *anaconda.TwitterApi

func GetApi() *anaconda.TwitterApi {
	if api != nil {
		return api
	}
	c := config.Env()
	anaconda.SetConsumerKey(c.ConsumerKey)
	anaconda.SetConsumerSecret(c.ConsumerSecret)
	api = anaconda.NewTwitterApi(c.AccessToken, c.AccessTokenSecret)
	return GetApi()
}

func GetAuthApi(u *account.User) *anaconda.TwitterApi {
	c := config.Env()
	anaconda.SetConsumerKey(c.ConsumerKey)
	anaconda.SetConsumerSecret(c.ConsumerSecret)
	return anaconda.NewTwitterApi(u.TwitterToken, u.TwitterSecret)
}

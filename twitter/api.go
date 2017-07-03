package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/remeh/smartwitter/config"
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

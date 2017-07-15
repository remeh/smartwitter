package twitter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/config"
	"github.com/remeh/smartwitter/log"
	. "github.com/remeh/smartwitter/twitter"

	"github.com/mrjones/oauth"
)

var tokens map[string]*oauth.RequestToken = make(map[string]*oauth.RequestToken)

type RedirectUserToTwitter struct{}

func (c RedirectUserToTwitter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cons := newConsumer()
	// TODO(remy): real route.
	tokenUrl := fmt.Sprintf("http://%s:9999/maketoken", r.Host)
	token, requestUrl, err := cons.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
		// TODO(remy):
		log.Error(err)
		w.WriteHeader(500)
		return
	}
	// Make sure to save the token, we'll need it for AuthorizeToken()
	tokens[token.Token] = token

	go func() {
		time.Sleep(60 * time.Second)
		delete(tokens, token.Token)
	}()

	http.Redirect(w, r, requestUrl, http.StatusTemporaryRedirect)
}

type GetTwitterToken struct{}

func (c GetTwitterToken) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cons := newConsumer()
	values := r.URL.Query()
	verificationCode := values.Get("oauth_verifier")
	tokenKey := values.Get("oauth_token")

	if _, ok := tokens[tokenKey]; !ok {
		// TODO(remy):!
		w.WriteHeader(400)
		return
	}

	atoken, err := cons.AuthorizeToken(tokens[tokenKey], verificationCode)
	defer delete(tokens, tokenKey)

	if err != nil {
		// TODO(remy):!
		w.WriteHeader(500)
		log.Error(err)
		return
	}

	api := GetAuthApi(atoken)
	tu, err := api.GetSelf(nil)
	if err != nil {
		w.WriteHeader(500)
		// TODO(remy):!
		log.Error(err)
		return
	}

	uid := account.GenTwitterUid(tu.IdStr)

	now := time.Now()

	u := &account.User{
		Uid: uid,

		CreationTime: now,
		LastLogin:    now,

		TwitterToken:  atoken.Token,
		TwitterSecret: atoken.Secret,

		TwitterId:   tu.IdStr,
		TwitterName: tu.Name,
	}

	if err := account.UserDAO().UpsertOnLogin(u); err != nil {
		w.WriteHeader(500)
		log.Error(err)
		return
	}

	// go to the application

	// TODO(remy): here we have the access token for this user
	// we want to get its IDs
}

func newConsumer() *oauth.Consumer {
	c := oauth.NewConsumer(
		config.Env().ConsumerKey,
		config.Env().ConsumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
	)
	c.Debug(true)
	return c
}

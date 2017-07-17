package twitter

import (
	"net/http"
	"time"

	"github.com/remeh/smartwitter/account"
	"github.com/remeh/smartwitter/api"
	"github.com/remeh/smartwitter/config"
	"github.com/remeh/smartwitter/log"
	. "github.com/remeh/smartwitter/twitter"

	"github.com/mrjones/oauth"
)

var tokens map[string]*oauth.RequestToken = make(map[string]*oauth.RequestToken)

type RedirectUserToTwitter struct{}

func (c RedirectUserToTwitter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cons := newConsumer()
	token, requestUrl, err := cons.GetRequestTokenAndUrl("")
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

	aapi := GetAuthApi(&account.User{
		TwitterToken:  atoken.Token,
		TwitterSecret: atoken.Secret,
	})
	tu, err := aapi.GetSelf(nil)
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

		TwitterToken:    atoken.Token,
		TwitterSecret:   atoken.Secret,
		TwitterId:       tu.IdStr,
		TwitterName:     tu.Name,
		TwitterUsername: tu.ScreenName,

		SessionToken: account.RandTok(),
	}

	if err := account.UserDAO().UpsertOnLogin(u); err != nil {
		// TODO(remy):!
		w.WriteHeader(500)
		log.Error(err)
		return
	}

	// set the session cookie and go to the application
	api.SetSessionCookie(w, u.SessionToken)

	http.Redirect(w, r, "https://twitter.remy.io", http.StatusTemporaryRedirect)
}

// ----------------------

type DebugSignIn struct{}

func (c DebugSignIn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conf := config.Env()
	if !conf.Debug {
		api.RenderBaseJson(w, 403, "route can't be called if not in debug mode")
		return
	}

	aapi := GetApi()
	tu, err := aapi.GetSelf(nil)
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	uid := account.GenTwitterUid(tu.IdStr)

	now := time.Now()

	tok, secr := conf.AccessToken, conf.AccessTokenSecret

	u := &account.User{
		Uid: uid,

		CreationTime: now,
		LastLogin:    now,

		TwitterToken:    tok,
		TwitterSecret:   secr,
		TwitterId:       tu.IdStr,
		TwitterName:     tu.Name,
		TwitterUsername: tu.ScreenName,

		SessionToken: account.RandTok(),
	}

	if err := account.UserDAO().UpsertOnLogin(u); err != nil {
		// TODO(remy):!
		w.WriteHeader(500)
		log.Error(err)
		return
	}

	// set the session cookie and go to the application
	api.SetSessionCookie(w, u.SessionToken)

	http.Redirect(w, r, "http://localhost:3000", http.StatusTemporaryRedirect)
}

// ----------------------

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

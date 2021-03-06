package config

import (
	"os"

	"github.com/remeh/smartwitter/log"
)

type Config struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string

	Conn       string
	AppUrl     string
	ListenAddr string
	PublicDir  string

	Debug bool

	loaded bool
}

var c Config

// EnvConfig returns the config set in
// environment variables.
// Could end the process if any of the env variable
// is unavailable.
func Env() Config {
	if c.loaded {
		return c
	}

	debug := os.Getenv("DEBUG") == "1"
	if debug {
		log.Warning("DEBUG DEBUG DEBUG mode !")
	}

	c = Config{
		ConsumerKey:       readEnvVar("CONSUMER_KEY", true, ""),
		ConsumerSecret:    readEnvVar("CONSUMER_SECRET", true, ""),
		AccessToken:       readEnvVar("ACCESS_TOKEN", false, ""),
		AccessTokenSecret: readEnvVar("ACCESS_TOKEN_SECRET", false, ""),
		Conn:              readEnvVar("CONN", false, "host=/var/run/postgresql sslmode=disable user=smartwitter dbname=smartwitter password=smartwitter"),
		AppUrl:            readEnvVar("APP_URL", false, "http://localhost:3000"),
		ListenAddr:        readEnvVar("ADDR", false, ":9999"),
		PublicDir:         readEnvVar("PUBLIC", false, "public/"),
		Debug:             debug,
		loaded:            true,
	}

	return c
}

func Reload() {
	c.loaded = false
}

func readEnvVar(v string, mandatory bool, def string) string {
	var rv string
	if rv = os.Getenv(v); len(rv) == 0 {
		if mandatory {
			log.Error("Can't find the environment variable:", v)
			os.Exit(1)
		} else {
			log.Warning("Using default value for:", v)
			rv = def
		}
	}
	return rv
}

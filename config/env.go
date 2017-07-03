package config

import (
	"log"
	"os"
)

const (
	ConsumerKey       string = "CONSUMER_KEY"
	ConsumerSecret    string = "CONSUMER_SECRET"
	AccessToken       string = "ACCESS_TOKEN"
	AccessTokenSecret string = "ACCESS_TOKEN_SECRET"
)

type Config struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// EnvConfig returns the config set in
// environment variables.
// Could end the process if any of the env variable
// is unavailable.
func Env() Config {
	return Config{
		ConsumerKey:       readEnvVar("CONSUMER_KEY"),
		ConsumerSecret:    readEnvVar("CONSUMER_SECRET"),
		AccessToken:       readEnvVar("ACCESS_TOKEN"),
		AccessTokenSecret: readEnvVar("ACCESS_TOKEN_SECRET"),
	}
}

func readEnvVar(v string) string {
	var rv string
	if rv = os.Getenv(v); len(rv) == 0 {
		log.Fatalln("Can't find the environment variable:", v)
	}
	return rv
}

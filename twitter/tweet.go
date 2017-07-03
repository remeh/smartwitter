package tweet

type Tweets []Tweet

type Tweet struct {
	Id     int64
	Time   time.Time
	Text   string
	Author TweetUser
}

type TweetUser struct {
	Description string
	ScreenName  string
	Name        string
	TimeZone    string
}

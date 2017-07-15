package account

import (
	"strings"
	"time"

	"github.com/remeh/uuid"

	"golang.org/x/crypto/bcrypt"
)

var (
	TwitterUidSeed = uuid.UUID("77e2c4cd-3fc8-4fbe-8ec8-b277dd1b5341")
)

type User struct {
	Uid uuid.UUID `json:"uid"`

	CreationTime time.Time
	LastLogin    time.Time

	TwitterToken    string
	TwitterSecret   string
	TwitterId       string
	TwitterUsername string
	TwitterName     string
}

// Crypt crypts the given password using bcrypt.
func Crypt(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(b), err
}

// Check validates that the hash is indeed derived from
// the given password.
func Check(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

// UnsubToken generates a random unsubscription
// token.
// It is composed of:
// - the char '1' (version)
// - first 8 chars of the user uid
// - with 3 randoms uuids (without -) appended.
func UnsubToken(uid uuid.UUID) string {
	end := randTok()
	start := uid.String()[0:8]
	return "1" + start + end
}

// Token generates a random password reset
// token.
// It is composed of:
// - the char '1' (version)
// - first 8 chars of the user uid
// - with 3 randoms uuids (without -) appended.
func PasswordResetToken(uid uuid.UUID) string {
	end := randTok()
	start := uid.String()[0:8]
	return "1" + start + end
}

func GenTwitterUid(twitterId string) uuid.UUID {
	return uuid.NewSHA1(TwitterUidSeed, []byte(twitterId))
}

// ----------------------

// randTok generates a random token composed of 3 uuids
// merge without the - char.
func randTok() string {
	rv := ""
	for i := 0; i < 3; i++ {
		str := uuid.New().String()
		str = strings.Replace(str, "-", "", -1)
		rv += str
	}
	return rv
}

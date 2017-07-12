package accounts

import (
	"strings"

	"github.com/remeh/uuid"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uid uuid.UUID `json:"uid"`

	UnsubToken  string `json:"-"`
	StripeToken string `json:"-"`
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

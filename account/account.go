package account

import (
	"github.com/remeh/uuid"
)

type Account struct {
	Uid       uuid.UUID
	Lastname  string
	Firstname string
	Email     string

	// ----------------------

	Usernames []string
	Terms     []string
}

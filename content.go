package main

import (
	"time"

	"github.com/remeh/uuid"
)

type Content struct {
	Id       uuid.UUID
	Author   string
	Text     string
	Hashtags []string
	Time     time.Time
}

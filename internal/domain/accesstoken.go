package domain

import (
	"time"
)

type AccessToken struct {
	ID             int
	UserID         int
	Token          string
	ExpirationTime time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

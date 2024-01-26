package domain

import "time"

type Wallet struct {
	UserID  int
	Balance float64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

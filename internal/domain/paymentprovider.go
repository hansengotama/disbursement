package domain

import (
	"github.com/google/uuid"
	"time"
)

type PaymentProvider struct {
	GUID     uuid.UUID
	Name     string
	AdminFee float64
	Type     string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

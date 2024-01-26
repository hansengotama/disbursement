package domain

import (
	"github.com/google/uuid"
	"time"
)

type DisbursementAccount struct {
	GUID                uuid.UUID
	UserID              int
	PaymentProviderGUID uuid.UUID
	Name                string
	Number              string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

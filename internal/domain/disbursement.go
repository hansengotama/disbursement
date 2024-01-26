package domain

import (
	"github.com/google/uuid"
	"time"
)

type Disbursement struct {
	ID                      int
	UserID                  int
	DisbursementAccountGUID uuid.UUID
	PaymentProviderGUID     uuid.UUID
	AccountName             string
	AccountNumber           string
	AdminFee                float64
	Amount                  float64
	AmountWithFee           float64
	Status                  string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

const (
	DisbursementStatusPending   = "pending"
	DisbursementStatusProcessed = "processed"
	DisbursementStatusFailed    = "failed"
)

package disbursementrepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
	"time"
)

type InsertDisbursementParam struct {
	Context                 context.Context
	Executor                postgres.SQLExecutor
	UserID                  int
	DisbursementAccountGUID uuid.UUID
	PaymentProviderGUID     uuid.UUID
	AccountName             string
	AccountNumber           string
	AdminFee                float64
	Amount                  float64
	AmountWithFee           float64
	Status                  string
}

type IDisbursementRepository interface {
	Insert(param InsertDisbursementParam) error
}

type DisbursementDB struct{}

func (r DisbursementDB) Insert(param InsertDisbursementParam) error {
	ctx, cancel := context.WithTimeout(param.Context, 2*time.Minute)
	defer cancel()

	_, err := param.Executor.ExecContext(ctx, "INSERT INTO disbursements(user_id, disbursement_account_guid, payment_provider_guid, account_name, account_number, admin_fee, amount, amount_with_fee, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", param.UserID, param.DisbursementAccountGUID, param.PaymentProviderGUID, param.AccountName, param.AccountNumber, param.AdminFee, param.Amount, param.AmountWithFee, param.Status)
	if err != nil {
		// logging
		return err
	}

	return nil
}

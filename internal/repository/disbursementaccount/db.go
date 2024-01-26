package disbursementaccountrepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
)

type InsertDisbursementAccountParam struct {
	Context             context.Context
	Executor            postgres.SQLExecutor
	GUID                uuid.UUID
	UserID              int
	PaymentProviderGUID uuid.UUID
	Name                string
	Number              string
}

type GetByGUIDParam struct {
	Context  context.Context
	Executor postgres.SQLExecutor
	GUID     uuid.UUID
}

type GetByGUIDParamRes struct {
	PaymentProviderGUID uuid.UUID
	Name                string
	Number              string
}

type IDisbursementAccountRepository interface {
	Insert(param InsertDisbursementAccountParam) error
	GetByGUID(param GetByGUIDParam) (*GetByGUIDParamRes, error)
}

type DisbursementAccountDB struct{}

func (r DisbursementAccountDB) Insert(param InsertDisbursementAccountParam) error {
	_, err := param.Executor.ExecContext(param.Context, "INSERT INTO disbursement_accounts(guid, user_id, payment_provider_guid, name, number) VALUES ($1, $2, $3, $4, $5)", param.GUID, param.UserID, param.PaymentProviderGUID, param.Name, param.Number)
	if err != nil {
		// logging
		return err
	}

	return nil
}

func (r DisbursementAccountDB) GetByGUID(param GetByGUIDParam) (*GetByGUIDParamRes, error) {
	row := param.Executor.QueryRowContext(param.Context, "SELECT payment_provider_guid, name, number from disbursement_accounts WHERE guid = $1", param.GUID)
	if row.Err() != nil {
		// logging
		return nil, row.Err()
	}

	var disbursementAccount GetByGUIDParamRes
	err := row.Scan(&disbursementAccount.PaymentProviderGUID, &disbursementAccount.Name, &disbursementAccount.Number)
	isNotFound := errors.Is(err, sql.ErrNoRows)
	if isNotFound {
		return nil, err
	}

	return &disbursementAccount, nil
}

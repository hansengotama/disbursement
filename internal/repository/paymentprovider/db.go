package paymentproviderrepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
	"time"
)

type InsertPaymentProviderParam struct {
	Context  context.Context
	Executor postgres.SQLExecutor
	GUID     uuid.UUID
	Name     string
	AdminFee float64
	Type     string
}

type GetAdminFeeByGUIDParam struct {
	Context  context.Context
	Executor postgres.SQLExecutor
	GUID     uuid.UUID
}

type IPaymentProviderRepository interface {
	Insert(param InsertPaymentProviderParam) error
	GetAdminFeeByGUID(param GetAdminFeeByGUIDParam) (float64, error)
}

type PaymentProviderDB struct{}

func (r PaymentProviderDB) Insert(param InsertPaymentProviderParam) error {
	ctx, cancel := context.WithTimeout(param.Context, 2*time.Minute)
	defer cancel()

	_, err := param.Executor.ExecContext(ctx, "INSERT INTO payment_providers(guid, name, admin_fee, type) VALUES ($1, $2, $3, $4)", param.GUID, param.Name, param.AdminFee, param.Type)
	if err != nil {
		// logging
		return err
	}

	return nil
}

func (r PaymentProviderDB) GetAdminFeeByGUID(param GetAdminFeeByGUIDParam) (float64, error) {
	ctx, cancel := context.WithTimeout(param.Context, 2*time.Minute)
	defer cancel()

	row := param.Executor.QueryRowContext(ctx, "SELECT admin_fee from payment_providers WHERE guid = $1", param.GUID)
	if row.Err() != nil {
		// logging
		return 0, row.Err()
	}

	var adminFee float64
	err := row.Scan(&adminFee)
	isNotFound := errors.Is(err, sql.ErrNoRows)
	if isNotFound {
		return 0, err
	}

	return adminFee, nil
}

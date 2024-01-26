package walletrepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
	"time"
)

type InsertWalletParam struct {
	Context  context.Context
	Executor postgres.SQLExecutor
	UserID   int
	Balance  float64
}

type GetWalletBalanceByUserIDParam struct {
	Context  context.Context
	Executor postgres.SQLExecutor
	UserID   int
}

type UpdateWalletBalanceByUserIDParam struct {
	Context  context.Context
	Executor postgres.SQLExecutor
	UserID   int
	Balance  float64
}

type IWalletRepository interface {
	Insert(param InsertWalletParam) error
	GetWalletBalanceByUserID(param GetWalletBalanceByUserIDParam) (balance float64, err error)
	UpdateWalletBalanceByUserID(param UpdateWalletBalanceByUserIDParam) error
}

type WalletDB struct{}

func (r WalletDB) Insert(param InsertWalletParam) error {
	_, err := param.Executor.ExecContext(param.Context, "INSERT INTO wallets(user_id, balance) VALUES ($1, $2)", param.UserID, param.Balance)
	if err != nil {
		// logging
		return err
	}

	return nil
}

func (r WalletDB) GetWalletBalanceByUserID(param GetWalletBalanceByUserIDParam) (float64, error) {
	row := param.Executor.QueryRowContext(param.Context, "SELECT balance from wallets WHERE user_id = $1", param.UserID)
	if row.Err() != nil {
		// logging
		return 0, row.Err()
	}

	var balance float64
	err := row.Scan(&balance)
	isNotFound := errors.Is(err, sql.ErrNoRows)
	if isNotFound {
		return 0, err
	}

	return balance, nil
}

func (r WalletDB) UpdateWalletBalanceByUserID(param UpdateWalletBalanceByUserIDParam) error {
	_, err := param.Executor.ExecContext(param.Context, "UPDATE wallets SET balance = $1, updated_at = $2 WHERE user_id = $3", param.Balance, time.Now(), param.UserID)
	if err != nil {
		// logging
		return err
	}

	return nil
}

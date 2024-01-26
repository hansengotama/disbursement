package userrepo

import (
	"context"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
	"time"
)

type InsertUserParam struct {
	Context  context.Context
	Executor postgres.SQLExecutor
	Name     string
}

type IUserRepository interface {
	Insert(param InsertUserParam) error
}

type UserDB struct{}

func (r UserDB) Insert(param InsertUserParam) error {
	ctx, cancel := context.WithTimeout(param.Context, 2*time.Minute)
	defer cancel()

	_, err := param.Executor.ExecContext(ctx, "INSERT INTO users(name) VALUES ($1)", param.Name)
	if err != nil {
		// logging
		return err
	}

	return nil
}

package userrepo

import (
	"context"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
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
	_, err := param.Executor.ExecContext(param.Context, "INSERT INTO users(name) VALUES ($1)", param.Name)
	if err != nil {
		// logging
		return err
	}

	return nil
}

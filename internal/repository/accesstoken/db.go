package accesstokenrepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
	"time"
)

type InsertAccessTokenParam struct {
	Context        context.Context
	Executor       postgres.SQLExecutor
	Token          string
	UserID         int
	ExpirationTime time.Time
}

type GetAccessTokenParam struct {
	Context  context.Context
	Executor postgres.SQLExecutor
	Token    string
}

type GetAccessTokenResponse struct {
	Token          string
	UserID         int
	ExpirationTime time.Time
}

type IAccessTokenRepository interface {
	Insert(param InsertAccessTokenParam) error
	GetAccessToken(param GetAccessTokenParam) (*GetAccessTokenResponse, error)
}

type AccessTokenDB struct{}

func (r AccessTokenDB) Insert(param InsertAccessTokenParam) error {
	_, err := param.Executor.ExecContext(param.Context, "INSERT INTO access_tokens(token, user_id, expiration_time) VALUES ($1, $2, $3)", param.Token, param.UserID, param.ExpirationTime)
	if err != nil {
		// logging
		return err
	}

	return nil
}

func (r AccessTokenDB) GetAccessToken(param GetAccessTokenParam) (*GetAccessTokenResponse, error) {
	row := param.Executor.QueryRowContext(param.Context, "SELECT token, user_id, expiration_time from access_tokens WHERE token = $1", param.Token)
	if row.Err() != nil {
		// logging
		return nil, row.Err()
	}

	var res GetAccessTokenResponse
	err := row.Scan(&res.Token, &res.UserID, &res.ExpirationTime)
	isNotFound := errors.Is(err, sql.ErrNoRows)
	if isNotFound {
		return nil, err
	}

	return &res, err
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hansengotama/disbursement/internal/lib/env"
	_ "github.com/lib/pq"
)

type SQLExecutor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

var dbConn *sql.DB

func init() {
	connectString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.GetPostgresDBUser(),
		env.GetPostgresDBPassword(),
		env.GetPostgresDBHost(),
		env.GetPostgresDBPort(),
		env.GetPostgresDBName(),
	)
	postgresDBCon, err := sql.Open("postgres", connectString)
	if err != nil {
		// logging
		panic(err)
	}

	dbConn = postgresDBCon
}

func GetConnection() *sql.DB {
	return dbConn
}

func CloseConnection() {
	err := dbConn.Close()
	if err != nil {
		// logging
		return
	}
}

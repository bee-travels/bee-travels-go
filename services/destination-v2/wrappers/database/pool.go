package database

import (
	"context"
	"github.com/bee-travels/bee-travels-go/destination-v2/wrappers/pgxpool"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type queryHook func(conn pgx.Conn, ctx context.Context, sql string, args ...interface{}) error

type Pool interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgxpool.Tx, error)
}

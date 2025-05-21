package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	Connection interface {
		Close()
		Begin(ctx context.Context) (pgx.Tx, error)
		Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
		Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
		QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
	}
	CatRepo struct {
		db Connection
	}
)

func NewCatRepo(c Connection) CatRepo {
	return CatRepo{
		db: c,
	}
}

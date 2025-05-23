package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/failuretoload/catdata/cat/domain"
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

func (r CatRepo) Query(ctx context.Context, input domain.QueryInput) ([]domain.CatRecord, error) {
	query := fmt.Sprintf("SELECT * FROM stats LIMIT %d OFFSET %d",
		input.Limit,
		input.Offset,
	)
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var records []domain.CatRecord
	for rows.Next() {
		var record domain.CatRecord
		var timestamp time.Time
		if err := rows.Scan(&record.ID, &record.Cat, &record.Weight, &record.Notes, &timestamp); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		record.Timestamp = timestamp.Format("01/02/2006 3:04 PM")
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return records, nil
}

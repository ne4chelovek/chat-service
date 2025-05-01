package log

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/ne4chelovek/chat_common/pkg/db"
	"github.com/ne4chelovek/chat_service/internal/repository"
)

const (
	tableLog = "transaction_log"
	columLog = "log"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.LogRepository {
	return &repo{db: db}
}

func (r *repo) Log(ctx context.Context, log string) error {
	builderInsert := sq.Insert(tableLog).
		PlaceholderFormat(sq.Dollar).
		Columns(columLog).
		Values(log).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "log_repository.Log",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

package balance

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/postikus/go-starter/model"
	"go.uber.org/zap"
	// std
	// 3rd
	// internal
)

type (
	executor interface {
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
)

type Repository struct {
	log  *zap.Logger
	conn executor
}

func NewRepository(log *zap.Logger, conn executor) *Repository {
	return &Repository{log: log.Named("repository.balance"), conn: conn}
}

func (r Repository) WithTx(exec executor) *Repository {
	r.conn = exec

	return &r
}

func (r *Repository) Insert(ctx context.Context, balance *model.Balance) (out *model.Balance, err error) {
	query, args, err := goqu.Dialect("mysql").
		Insert("balance").
		Rows(goqu.Record{
			"user_id": balance.UserID,
			"amount":  balance.Amount,
		}).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	res, err := r.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute insert statement: %w", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected count: %w", err)
	}
	if ra < 1 {
		return nil, fmt.Errorf("no rows affected")
	}

	return balance, nil
}

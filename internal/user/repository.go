package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/postikus/go-starter/model"
	"go.uber.org/zap"
	// std
	// 3rd
	// internal
)

type (
	beginner interface {
		BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	}
	executor interface {
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
	committer interface {
		Close(tx *sqlx.Tx) func(error)
	}
)

type Repository struct {
	log       *zap.Logger
	conn      executor
	beginner  beginner
	committer committer
}

func NewRepository(log *zap.Logger, conn executor, beginner beginner, committer committer) *Repository {
	return &Repository{
		log:       log.Named("repository.user"),
		conn:      conn,
		beginner:  beginner,
		committer: committer,
	}
}

func (r Repository) WithTx(exec executor) *Repository {
	r.conn = exec

	return &r
}

func (r Repository) StartTx(ctx context.Context) (*sqlx.Tx, func(error), error) {
	tx, err := r.beginner.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	return tx, r.committer.Close(tx), nil
}

func (r *Repository) Insert(ctx context.Context, user *model.User) (out *model.User, err error) {
	query, args, err := goqu.Dialect("mysql").
		Insert("reg_user").
		Rows(goqu.Record{
			"name": user.Name,
		}).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	res, err := r.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute insert statement: %w", err)
	}

	user.ID, err = res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected count: %w", err)
	}

	return user, nil
}

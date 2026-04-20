package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("postgres: not found")

type Conn interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type TxOpener interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Repository struct {
	conn   Conn
	opener TxOpener
}

type TxRepository struct {
	tx *sql.Tx
	*Repository
}

func NewRepository(conn *sql.DB) *Repository {
	return &Repository{
		conn:   conn,
		opener: conn,
	}
}

func (x *TxRepository) Commit(_ context.Context) error {
	return x.tx.Commit()
}

func (x *TxRepository) Rollback(_ context.Context) error {
	return x.tx.Rollback()
}

func (r *Repository) OpenTx(ctx context.Context) (*TxRepository, error) {
	tx, err := r.opener.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &TxRepository{
		tx: tx,
		Repository: &Repository{
			conn: tx,
		},
	}, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func mapNotFound(entity string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s: %w", entity, ErrNotFound)
	}

	return err
}

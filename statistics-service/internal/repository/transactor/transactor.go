package transactor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Queries interface {
	Exec(query string, args ...any) (commandTag sql.Result, err error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type txKey struct{}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func GetQueries(ctx context.Context, defaultQueries Queries) Queries {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}

	return defaultQueries
}

type TxBeginner struct {
	db *sql.DB
}

func NewTxBeginner(db *sql.DB) *TxBeginner {
	return &TxBeginner{db: db}
}

func (r *TxBeginner) WithTransactionValue(ctx context.Context, txFunc func(ctx context.Context) (any, error)) (any, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		}
	}()

	result, err := txFunc(injectTx(ctx, tx))
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction error: %w", err)
	}

	return result, nil
}

func (r *TxBeginner) WithTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		}
	}()

	err = txFunc(injectTx(ctx, tx))
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction error: %w", err)
	}

	return nil
}

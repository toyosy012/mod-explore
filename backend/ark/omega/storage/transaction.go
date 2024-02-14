package storage

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

func (c *Client[T, ID]) WithTransaction(ctx context.Context, fn func(context.Context) (any, error)) (_ any, err error) {
	timeout, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	tx, err := c.BeginTxx(timeout, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		return nil, err
	}

	defer func(err error) {
		if p := recover(); p != nil {
			if err = tx.Rollback(); err != nil {
				c.logger.ErrorContext(ctx, "failed to rollback transaction in panic", slog.Any("error", err))
			}
			panic(p) // rethrow
		}
		if err != nil {
			if err = tx.Rollback(); err != nil {
				c.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("error", err))
			}
			return
		}
		if err = tx.Commit(); err != nil {
			c.logger.ErrorContext(ctx, "failed to commit transaction", slog.Any("error", err))
		}
	}(err)

	return fn(SetTx(timeout, tx))
}

type txKey struct{}

func SetTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	_, already := ctx.Value(txKey{}).(*sqlx.Tx)
	if already {
		return ctx
	}

	return context.WithValue(ctx, txKey{}, tx)
}

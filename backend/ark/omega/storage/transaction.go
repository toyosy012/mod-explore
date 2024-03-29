package storage

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"

	"mods-explore/ark/omega/logic/variant/domain/service"
)

func (c *Client) WithTransaction(ctx context.Context, fn func(context.Context) (any, error)) (_ any, err error) {
	timeout, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	tx, err := c.BeginTxx(timeout, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			if e := tx.Rollback(); e != nil {
				c.logger.ErrorContext(ctx, "failed to rollback transaction in panic", slog.Any("error", e))
			}
			err = service.IntervalServerError
			return
		}
		if err != nil {
			if e := tx.Rollback(); e != nil {
				c.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("error", e))
			}
			return
		}
		if e := tx.Commit(); e != nil {
			c.logger.ErrorContext(ctx, "failed to commit transaction", slog.Any("error", e))
		}
	}()
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

package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"mods-explore/ark/omega/logic/variant/domain/service"
)

// Client sqlxのインスタンスは使いまわしたいのでテーブル毎にクライアントモジュールを生成できるようにする
// TODO Tをanyにせず、IDの型パラメータを指定しなくてもいいようにしたい。もしくはIDをテーブルのidの型に強制的に一致するようにしたい。
type Client[T any, ID any] struct {
	*sqlx.DB
	logger *slog.Logger
}

func Connect[T any, ID any](db *sqlx.DB, logger *slog.Logger) (_ *Client[T, ID], err error) {
	if logger == nil {
		logger = slog.Default()
	}

	return &Client[T, ID]{
		DB:     db,
		logger: logger,
	}, nil
}

func (c *Client[T, ID]) NamedGet(ctx context.Context, query string, args ...any) (*T, error) {
	query, args, err := c.BindNamed(query, args)
	if err != nil {
		return nil, err
	}

	var row T
	if err = c.GetContext(ctx, &row, query, args...); errors.Is(err, sql.ErrNoRows) {
		return nil, service.NotFound
	} else if err != nil {
		return nil, err
	}

	return &row, nil
}

func (c *Client[T, ID]) NamedSelect(ctx context.Context, query string) ([]T, error) {
	var rows []T
	if err := c.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	return rows, nil
}

func (c *Client[T, ID]) NamedStore(ctx context.Context, query string, arg any) (id ID, err error) {
	stmt, err := c.PrepareNamedContext(ctx, query)
	if err != nil {
		return id, err
	}

	err = stmt.QueryRowContext(ctx, arg).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (c *Client[T, ID]) NamedDelete(ctx context.Context, query string, arg any) error {
	_, err := c.NamedExecContext(
		ctx,
		query,
		arg,
	)
	return err
}

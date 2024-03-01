package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/samber/do"

	"mods-explore/ark/omega/logic/variant/domain/service"
)

// Client sqlxのインスタンスは使いまわしたいのでテーブル毎にクライアントモジュールを生成できるようにする
// TODO Tをanyにせず、IDの型パラメータを指定しなくてもいいようにしたい。もしくはIDをテーブルのidの型に強制的に一致するようにしたい。
type Client struct {
	*sqlx.DB
	logger *slog.Logger
}

func NewSQLxClient(injector *do.Injector) (*Client, error) {
	return &Client{
		DB:     do.MustInvoke[*sqlx.DB](injector),
		logger: do.MustInvoke[*slog.Logger](injector),
	}, nil
}

func NamedGet[T any](ctx context.Context, c *Client, query string, args ...any) (*T, error) {
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

func Select[T any](ctx context.Context, c *Client, query string) ([]T, error) {
	var rows []T
	if err := c.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	return rows, nil
}

func NamedStore[ID any](ctx context.Context, c *Client, query string, arg any) (id ID, err error) {
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

func NamedDelete(ctx context.Context, c *Client, query string, arg any) error {
	_, err := c.NamedExecContext(
		ctx,
		query,
		arg,
	)
	return err
}

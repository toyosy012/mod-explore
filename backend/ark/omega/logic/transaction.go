package logic

import "context"

type Transactioner interface {
	WithTransaction(ctx context.Context, fn func(context.Context) (any, error)) (any, error)
}

type txerKey struct{}

func SetTransctioner(ctx context.Context, tx Transactioner) context.Context {
	return context.WithValue(ctx, txerKey{}, tx)
}

func GetTransctioner(ctx context.Context) (Transactioner, bool) {
	v := ctx.Value(txerKey{})
	if v != nil {
		return v.(Transactioner), true
	}
	return nil, false
}

// UseTransactioner 戻り値にデータもあるINSERTなどの処理に用いる
func UseTransactioner[R any](ctx context.Context, fn func(ctx context.Context) (R, error)) (_ R, err error) {
	tx, ok := GetTransctioner(ctx)
	if !ok {
		return fn(ctx)
	}
	res, err := tx.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		return fn(ctx)
	})
	if err != nil {
		return
	}

	return res.(R), nil
}

// UseTransactioner0 戻り値がエラーのみのDELETEなどの処理に用いる
func UseTransactioner0(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, ok := GetTransctioner(ctx)
	if !ok {
		return fn(ctx)
	}
	_, err := tx.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		return nil, fn(ctx)
	})

	return err
}

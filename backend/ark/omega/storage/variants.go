package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"mods-explore/ark/omega/logic/variant/domain/model"
	"mods-explore/ark/omega/logic/variant/domain/service"
)

type VariantClient struct {
	*sqlx.DB
}

func NewVariantClient(dsn string) (service.VariantRepository, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err = db.PingContext(timeout); err != nil {
		return nil, fmt.Errorf("connection timeout: %s ", err)
	}

	return VariantClient{
		DB: db,
	}, nil
}

// variantModel Listでも1度に取得されるレコード量は決まっているので、domain modelで異なるバインド用モデルを定義する
type variantModel struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Group string `db:"group"`
}

func (v VariantClient) FindVariant(ctx context.Context, id model.VariantID) (*model.Variant, error) {
	query, args, err := v.BindNamed(
		`SELECT variants.id, variants.name, groups.name AS "group" FROM variants
    INNER JOIN groups ON (variants.group_id = groups.id) WHERE variants.id = :id;`,
		map[string]any{"id": id},
	)

	if err != nil {
		return nil, err
	}

	var row variantModel
	if err = v.GetContext(ctx, &row, query, args...); errors.Is(err, sql.ErrNoRows) {
		return nil, service.NotFound
	} else if err != nil {
		return nil, err
	}

	var variant = model.NewVariant(
		model.VariantID(row.ID),
		model.GroupName(row.Group),
		model.Name(row.Name),
	)
	return &variant, nil
}

func (v VariantClient) ListVariants(ctx context.Context) (model.Variants, error) {
	query := `SELECT variants.id, variants.name, groups.name AS "group" FROM variants
    INNER JOIN groups ON (variants.group_id = groups.id);`

	var rows []variantModel
	if err := v.SelectContext(ctx, &rows, query); errors.Is(err, sql.ErrNoRows) {
		return nil, service.NotFound
	} else if err != nil {
		return nil, err
	}

	var results model.Variants
	for _, r := range rows {
		results = append(
			results,
			model.NewVariant(
				model.VariantID(r.ID),
				model.GroupName(r.Group),
				model.Name(r.Name),
			),
		)
	}

	return results, nil
}

func (v VariantClient) CreateVariant(ctx context.Context, create service.CreateVariant) (*model.Variant, error) {
	stmt, err := v.PrepareNamedContext(
		ctx,
		`INSERT INTO variants (name, group_id) VALUES (:name, :groupID) RETURNING id;`,
	)
	if err != nil {
		return nil, err
	}

	var id int
	err = stmt.
		QueryRowxContext(ctx, map[string]any{"name": create.Name(), "groupID": create.GroupID()}).
		Scan(&id)
	if err != nil {
		return nil, err
	}

	result, err := v.FindVariant(ctx, model.VariantID(id))
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (v VariantClient) UpdateVariant(context.Context, service.UpdateVariant) (*model.Variant, error) {
	return nil, nil
}
func (v VariantClient) DeleteVariant(context.Context, model.VariantID) error { return nil }

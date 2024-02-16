package storage

import (
	"context"

	"github.com/samber/do"

	"mods-explore/ark/omega/logic/variant/domain/model"
	"mods-explore/ark/omega/logic/variant/domain/service"
)

type VariantClient struct {
	*Client[VariantModel, int]
}

func NewVariantClient(injector *do.Injector) (service.VariantRepository, error) {
	return VariantClient{
		do.MustInvoke[*Client[VariantModel, int]](injector),
	}, nil
}

// VariantModel Listでも1度に取得されるレコード量は決まっているので、domain modelで異なるバインド用モデルを定義する
type VariantModel struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Group string `db:"group"`
}

func (v VariantClient) FindVariant(ctx context.Context, id model.VariantID) (*model.Variant, error) {
	row, err := v.NamedGet(
		ctx,
		`SELECT variants.id, variants.name, groups.name AS "group" FROM variants
    INNER JOIN groups ON (variants.group_id = groups.id) WHERE variants.id = :id;`,
		map[string]any{"id": id},
	)
	if err != nil {
		return nil, err
	}

	var variant = model.NewVariant(
		model.VariantID(row.ID),
		model.VariantGroupName(row.Group),
		model.Name(row.Name),
	)
	return &variant, nil
}

func (v VariantClient) ListVariants(ctx context.Context) (model.Variants, error) {
	rows, err := v.NamedSelect(
		ctx,
		`SELECT variants.id, variants.name, groups.name AS "group" FROM variants
    INNER JOIN groups ON (variants.group_id = groups.id);`,
	)
	if err != nil {
		return nil, err
	}

	var results model.Variants
	for _, r := range rows {
		results = append(
			results,
			model.NewVariant(
				model.VariantID(r.ID),
				model.VariantGroupName(r.Group),
				model.Name(r.Name),
			),
		)
	}

	return results, nil
}

func (v VariantClient) CreateVariant(ctx context.Context, create service.CreateVariant) (*model.Variant, error) {
	id, err := v.NamedStore(
		ctx,
		`INSERT INTO variants (name, group_id) VALUES (:name, :groupID) RETURNING id;`,
		map[string]any{"name": create.Name(), "groupID": create.GroupID()},
	)
	if err != nil {
		return nil, err
	}

	result, err := v.FindVariant(ctx, model.VariantID(id))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (v VariantClient) UpdateVariant(ctx context.Context, update service.UpdateVariant) (*model.Variant, error) {
	id, err := v.NamedStore(
		ctx,
		`UPDATE variants SET name = :name, group_id = :groupID, updated_at = NOW() WHERE id = :id RETURNING id;`,
		map[string]any{"id": update.ID(), "name": update.Name(), "groupID": update.GroupID()},
	)
	if err != nil {
		return nil, err
	}

	result, err := v.FindVariant(ctx, model.VariantID(id))
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (v VariantClient) DeleteVariant(ctx context.Context, id model.VariantID) error {
	return v.NamedDelete(ctx, `DELETE FROM variants WHERE id = :id;`, map[string]any{"id": id})
}

package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log/slog"

	"mods-explore/ark/omega/logic/variant/domain/model"
	"mods-explore/ark/omega/logic/variant/domain/service"
)

// GroupVariantModel variantsに集約しても良さそうだったがgroups単体で取り扱う可能性があるので分離しておく
type variantGroupModel struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type VariantGroupClient struct {
	*Client[variantGroupModel, int]
}

func NewVariantGroupClient(db *sqlx.DB, logger *slog.Logger) (service.VariantGroupRepository, error) {
	cli, err := NewSQLxClient[variantGroupModel, int](db, logger)
	if err != nil {
		return nil, err
	}

	return VariantGroupClient{
		cli,
	}, nil
}

func (v VariantGroupClient) Select(ctx context.Context, id model.VariantGroupID) (*model.VariantGroup, error) {
	row, err := v.NamedGet(
		ctx,
		`SELECT id, name FROM groups WHERE id = :id;`,
		map[string]any{"id": id},
	)
	if err != nil {
		return nil, err
	}

	var variant = model.NewVariantGroup(
		model.VariantGroupID(row.ID),
		model.VariantGroupName(row.Name),
	)
	return &variant, nil
}

func (v VariantGroupClient) List(ctx context.Context) (model.VariantGroups, error) {
	rows, err := v.NamedSelect(
		ctx,
		`SELECT id, name FROM groups;`,
	)
	if err != nil {
		return nil, err
	}

	var results model.VariantGroups
	for _, r := range rows {
		results = append(
			results,
			model.NewVariantGroup(
				model.VariantGroupID(r.ID),
				model.VariantGroupName(r.Name),
			),
		)
	}

	return results, nil
}

func (v VariantGroupClient) Insert(ctx context.Context, create service.CreateVariantGroup) (*model.VariantGroup, error) {
	id, err := v.NamedStore(
		ctx,
		`INSERT INTO groups (name) VALUES (:name) RETURNING id;`,
		map[string]any{"name": create.Name()},
	)
	if err != nil {
		return nil, err
	}

	result, err := v.Select(ctx, model.VariantGroupID(id))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (v VariantGroupClient) Update(ctx context.Context, update service.UpdateVariantGroup) (*model.VariantGroup, error) {
	_, err := v.NamedStore(
		ctx,
		`UPDATE groups SET name = :name, updated_at = NOW() WHERE id = :id;`,
		map[string]any{"id": update.ID(), "name": update.Name()},
	)
	if err != nil {
		return nil, err
	}

	result, err := v.Select(ctx, update.ID())
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (v VariantGroupClient) Delete(ctx context.Context, id model.VariantGroupID) error {
	return v.NamedDelete(ctx, `DELETE FROM groups WHERE id = :id;`, map[string]any{"id": id})
}

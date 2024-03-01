package storage

import (
	"context"

	"github.com/samber/do"

	"mods-explore/ark/omega/logic/creature/domain/model"
	"mods-explore/ark/omega/logic/creature/domain/service"
)

type DinosaurModel struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	BaseHealth int    `db:"health"`
	BaseMelee  int    `db:"melee"`
}

type DinosaurClient struct {
	*Client
}

func NewDinosaurClient(injector *do.Injector) (service.DinosaurCommandRepository, error) {
	return DinosaurClient{
		do.MustInvoke[*Client](injector),
	}, nil
}

func (c DinosaurClient) Insert(ctx context.Context, create service.CreateDinosaur) (model.DinosaurID, error) {
	id, err := NamedStore[int](
		ctx,
		c.Client,
		`INSERT INTO dinosaurs (name, health, melee) VALUES (:name, :health, :melee) RETURNING id;`,
		map[string]any{"name": create.Name(), "health": create.Health(), "melee": create.Melee()},
	)
	if err != nil {
		return 0, err
	}
	return model.DinosaurID(id), err
}

func (c DinosaurClient) Update(ctx context.Context, update service.UpdateDinosaur) error {
	_, err := NamedStore[int](
		ctx,
		c.Client,
		`UPDATE dinosaurs SET name = :name, health = :health, melee = :melee, updated_at = NOW() WHERE id = :id;`,
		map[string]any{"id": update.ID(), "name": update.Name(), "health": update.Health(), "melee": update.Melee()},
	)
	return err
}
func (c DinosaurClient) Delete(ctx context.Context, id model.DinosaurID) error {
	return NamedDelete(ctx, c.Client, `DELETE FROM dinosaurs WHERE id = :id;`, map[string]any{"id": id})
}

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
	*Client[DinosaurModel, int]
}

func NewDinosaurClient(injector *do.Injector) (service.DinosaurCommandRepository, error) {
	return DinosaurClient{
		do.MustInvoke[*Client[DinosaurModel, int]](injector),
	}, nil
}

func (c DinosaurClient) Insert(ctx context.Context, create service.CreateDinosaur) error {
	_, err := c.NamedStore(
		ctx,
		`INSERT INTO dinosaurs (name, health, melee) VALUES (:name, :health, :melee);`,
		map[string]any{"name": create.Name(), "health": create.Health(), "melee": create.Melee()},
	)
	return err
}

func (c DinosaurClient) Update(ctx context.Context, update service.UpdateDinosaur) error {
	_, err := c.NamedStore(
		ctx,
		`UPDATE dinosaurs SET name = :name, health = :health, melee = :melee, updated_at = NOW() WHERE id = :id;`,
		map[string]any{"id": update.ID(), "name": update.Name(), "health": update.Health(), "melee": update.Melee()},
	)
	return err
}
func (c DinosaurClient) Delete(ctx context.Context, id model.DinosaurID) error {
	return c.NamedDelete(ctx, `DELETE FROM dinosaurs WHERE id = :id;`, map[string]any{"id": id})
}

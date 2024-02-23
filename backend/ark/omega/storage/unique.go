package storage

import (
	"context"

	"mods-explore/ark/omega/logic/creature/domain/model"
	"mods-explore/ark/omega/logic/creature/domain/service"
)

type UniqueModel struct {
	ID               int     `db:"id"`
	Name             string  `db:"name"`
	HealthMultiplier float32 `db:"health_multiplier"`
	DamageMultiplier float32 `db:"damage_multiplier"`
}

type UniqueCommandRepo struct {
	Client[UniqueModel, int]
}

func (r UniqueCommandRepo) Insert(ctx context.Context, create service.CreateUniqueDinosaur) (model.UniqueDinosaurID, error) {
	return model.UniqueDinosaurID(0), nil
}

func (r UniqueCommandRepo) Update(ctx context.Context, update service.UpdateUniqueDinosaur) error {
	return nil
}

func (r UniqueCommandRepo) Delete(ctx context.Context, id model.UniqueDinosaurID) error {
	return nil
}

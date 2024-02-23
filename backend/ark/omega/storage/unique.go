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
	id, err := r.NamedStore(
		ctx,
		`INSERT INTO uniques (dinosaur_id, name, health_multiplier, damage_multiplier)
			VALUES (:dinosaur_id, :name, :health_multiplier, :damage_multiplier)
			RETURNING id;`,
		map[string]any{
			"dinosaur_id": create.DinosaurID(), "name": create.Name(),
			"health_multiplier": create.HealthMultiplier(), "damage_multiplier": create.DamageMultiplier()},
	)
	if err != nil {
		return model.UniqueDinosaurID(0), err
	}
	return model.UniqueDinosaurID(id), nil
}

func (r UniqueCommandRepo) Update(ctx context.Context, update service.UpdateUniqueDinosaur) error {
	_, err := r.NamedStore(
		ctx,
		`UPDATE uniques 
			SET dinosaur_id = :dinosaur_id, name = :name, 
			    health_multiplier = :health_multiplier, damage_multiplier = :damage_multiplier, updated_at = NOW() 
			WHERE id = :id;`,
		map[string]any{
			"id": update.ID(), "dinosaur_id": update.DinosaurID(), "name": update.Name(),
			"health_multiplier": update.HealthMultiplier(), "damage_multiplier": update.DamageMultiplier(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r UniqueCommandRepo) Delete(ctx context.Context, id model.UniqueDinosaurID) error {
	return nil
}

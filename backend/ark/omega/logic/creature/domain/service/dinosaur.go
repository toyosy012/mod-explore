package service

import (
	"context"

	"mods-explore/ark/omega/logic/creature/domain/model"
)

type DinosaurCommandRepository interface {
	Insert(context.Context, CreateDinosaur) (*model.Dinosaur, error)
	Update(context.Context, UpdateDinosaur) (*model.Dinosaur, error)
	Delete(context.Context, model.UniqueDinosaurID) error
}

type CreateDinosaur struct {
	name       model.DinosaurName
	baseHealth model.Health
	baseMelee  model.Melee
}

func NewCreateDinosaur(
	name model.DinosaurName,
	baseHealth model.Health,
	baseMelee model.Melee,
) CreateDinosaur {
	return CreateDinosaur{
		name:       name,
		baseHealth: baseHealth,
		baseMelee:  baseMelee,
	}
}

type UpdateDinosaur struct {
	id         model.DinosaurID
	name       model.DinosaurName
	baseHealth model.Health
	baseMelee  model.Melee
}

func NewUpdateDinosaur(
	id model.DinosaurID,
	name model.DinosaurName,
	baseHealth model.Health,
	baseMelee model.Melee,
) UpdateDinosaur {
	return UpdateDinosaur{
		id:         id,
		name:       name,
		baseHealth: baseHealth,
		baseMelee:  baseMelee,
	}
}

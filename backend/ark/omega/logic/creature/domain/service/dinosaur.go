package service

import (
	"context"

	"mods-explore/ark/omega/logic/creature/domain/model"
)

type DinosaurCommandRepository interface {
	Insert(context.Context, CreateDinosaur) (*ResponseDinosaur, error)
	Update(context.Context, UpdateDinosaur) (*ResponseDinosaur, error)
	Delete(context.Context, model.DinosaurID) error
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

type ResponseDinosaur struct {
	id         model.DinosaurID
	name       model.DinosaurName
	baseHealth model.Health
	baseMelee  model.Melee
}

func NewResponseDinosaur(
	id model.DinosaurID,
	name model.DinosaurName,
	baseHealth model.Health,
	baseMelee model.Melee,
) ResponseDinosaur {
	return ResponseDinosaur{
		id:         id,
		name:       name,
		baseHealth: baseHealth,
		baseMelee:  baseMelee,
	}
}

func (d ResponseDinosaur) ID() model.DinosaurID     { return d.id }
func (d ResponseDinosaur) Name() model.DinosaurName { return d.name }
func (d ResponseDinosaur) Health() model.Health     { return d.baseHealth }
func (d ResponseDinosaur) Melee() model.Melee       { return d.baseMelee }

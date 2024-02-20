package service

import (
	"context"

	"mods-explore/ark/omega/logic/creature/domain/model"
)

type DinosaurCommandRepository interface {
	Insert(context.Context, CreateDinosaur) (model.DinosaurID, error)
	Update(context.Context, UpdateDinosaur) error
	Delete(context.Context, model.DinosaurID) error
}

type CreateDinosaur struct {
	name       model.DinosaurName
	baseHealth model.Health
	baseMelee  model.Melee
}

func (d CreateDinosaur) Name() model.DinosaurName { return d.name }
func (d CreateDinosaur) Health() model.Health     { return d.baseHealth }
func (d CreateDinosaur) Melee() model.Melee       { return d.baseMelee }

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

func (d UpdateDinosaur) ID() model.DinosaurID     { return d.id }
func (d UpdateDinosaur) Name() model.DinosaurName { return d.name }
func (d UpdateDinosaur) Health() model.Health     { return d.baseHealth }
func (d UpdateDinosaur) Melee() model.Melee       { return d.baseMelee }

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

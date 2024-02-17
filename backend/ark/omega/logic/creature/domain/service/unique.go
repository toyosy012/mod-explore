package service

import (
	"context"

	"mods-explore/ark/omega/logic/creature/domain/model"
)

type UniqueRepository interface {
	Select(context.Context, model.UniqueDinosaurID) (*model.UniqueDinosaur, error)
	List(context.Context) (model.UniqueDinosaurs, error)
	Insert(context.Context, CreateUniqueDinosaur) (*model.UniqueDinosaur, error)
	Update(context.Context, UpdateUniqueDinosaur) (*model.UniqueDinosaur, error)
	Delete(context.Context, model.UniqueDinosaurID) error
}

type CreateUniqueDinosaur struct {
	CreateDinosaur
	name             model.UniqueName
	variants         model.UniqueVariant
	healthMultiplier model.UniqueMultiplier[model.Health]
	damageMultiplier model.UniqueMultiplier[model.Melee]
}

func NewCreateUniqueDinosaur(
	base CreateDinosaur,
	name model.UniqueName,
	variants model.UniqueVariant,
	healthMultiplier model.UniqueMultiplier[model.Health],
	damageMultiplier model.UniqueMultiplier[model.Melee],
) CreateUniqueDinosaur {
	return CreateUniqueDinosaur{
		CreateDinosaur:   base,
		name:             name,
		variants:         variants,
		healthMultiplier: healthMultiplier,
		damageMultiplier: damageMultiplier,
	}
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

type UpdateUniqueDinosaur struct {
	UpdateDinosaur
	uniqueDinoID     model.UniqueDinosaurID
	name             model.UniqueName
	variants         model.UniqueVariant
	healthMultiplier model.UniqueMultiplier[model.Health]
	damageMultiplier model.UniqueMultiplier[model.Melee]
}

func NewUpdateUniqueDinosaur(
	base UpdateDinosaur,
	uniqueDinoID model.UniqueDinosaurID,
	name model.UniqueName,
	variants model.UniqueVariant,
	healthMultiplier model.UniqueMultiplier[model.Health],
	damageMultiplier model.UniqueMultiplier[model.Melee],
) UpdateUniqueDinosaur {
	return UpdateUniqueDinosaur{
		UpdateDinosaur:   base,
		uniqueDinoID:     uniqueDinoID,
		name:             name,
		variants:         variants,
		healthMultiplier: healthMultiplier,
		damageMultiplier: damageMultiplier,
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

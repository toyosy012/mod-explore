package service

import (
	"context"

	"mods-explore/ark/omega/logic/creature/domain/model"
)

// UniqueQueryRepository 集約内のテーブルをjoinしてレコードを取得する処理を定義
type UniqueQueryRepository interface {
	Select(context.Context, model.UniqueDinosaurID) (*model.UniqueDinosaur, error)
	List(context.Context) (model.UniqueDinosaurs, error)
}

type UniqueCommandRepository interface {
	Insert(context.Context, CreateUniqueDinosaur) (*ResponseUnique, error)
	Update(context.Context, UpdateUniqueDinosaur) (*ResponseUnique, error)
	Delete(context.Context, model.UniqueDinosaurID) error
}

type CreateCreature struct {
	CreateDinosaur
	CreateUniqueDinosaur
	CreateVariants
}

func NewCreateCreature(
	base CreateDinosaur,
	unique CreateUniqueDinosaur,
	variants CreateVariants,
) CreateCreature {
	return CreateCreature{
		CreateDinosaur:       base,
		CreateUniqueDinosaur: unique,
		CreateVariants:       variants,
	}
}

type CreateUniqueDinosaur struct {
	name             model.UniqueName
	healthMultiplier model.UniqueMultiplier[model.Health]
	damageMultiplier model.UniqueMultiplier[model.Melee]
}

func NewCreateUniqueDinosaur(
	name model.UniqueName,
	healthMultiplier model.UniqueMultiplier[model.Health],
	damageMultiplier model.UniqueMultiplier[model.Melee],
) CreateUniqueDinosaur {
	return CreateUniqueDinosaur{
		name:             name,
		healthMultiplier: healthMultiplier,
		damageMultiplier: damageMultiplier,
	}
}

type UpdateCreature struct {
	UpdateDinosaur
	UpdateUniqueDinosaur
	UpdateVariants
}

func NewUpdateCreature(
	base UpdateDinosaur,
	unique UpdateUniqueDinosaur,
	variants UpdateVariants,
) UpdateCreature {
	return UpdateCreature{
		UpdateDinosaur:       base,
		UpdateUniqueDinosaur: unique,
		UpdateVariants:       variants,
	}
}

type UpdateUniqueDinosaur struct {
	uniqueDinoID     model.UniqueDinosaurID
	name             model.UniqueName
	healthMultiplier model.UniqueMultiplier[model.Health]
	damageMultiplier model.UniqueMultiplier[model.Melee]
}

func NewUpdateUniqueDinosaur(
	uniqueDinoID model.UniqueDinosaurID,
	name model.UniqueName,
	healthMultiplier model.UniqueMultiplier[model.Health],
	damageMultiplier model.UniqueMultiplier[model.Melee],
) UpdateUniqueDinosaur {
	return UpdateUniqueDinosaur{
		uniqueDinoID:     uniqueDinoID,
		name:             name,
		healthMultiplier: healthMultiplier,
		damageMultiplier: damageMultiplier,
	}
}

func (d UpdateUniqueDinosaur) ID() model.UniqueDinosaurID { return d.uniqueDinoID }

type ResponseUnique struct {
	id               model.UniqueDinosaurID
	name             model.UniqueName
	healthMultiplier model.UniqueMultiplier[model.Health]
	damageMultiplier model.UniqueMultiplier[model.Melee]
}

func NewResponseUnique(
	id model.UniqueDinosaurID,
	name model.UniqueName,
	healthMultiplier model.UniqueMultiplier[model.Health],
	damageMultiplier model.UniqueMultiplier[model.Melee],
) ResponseUnique {
	return ResponseUnique{id, name, healthMultiplier, damageMultiplier}
}

func (u ResponseUnique) ID() model.UniqueDinosaurID { return u.id }
func (u ResponseUnique) Name() model.UniqueName     { return u.name }
func (u ResponseUnique) HealthMultiplier() model.UniqueMultiplier[model.Health] {
	return u.healthMultiplier
}
func (u ResponseUnique) MeleeMultiplier() model.UniqueMultiplier[model.Melee] {
	return u.damageMultiplier
}

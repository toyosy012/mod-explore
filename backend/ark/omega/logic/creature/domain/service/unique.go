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
	Insert(context.Context, CreateUniqueDinosaur) (model.UniqueDinosaurID, error)
	Update(context.Context, UpdateUniqueDinosaur) error
	Delete(context.Context, model.UniqueDinosaurID) error
}

type CreateCreature struct {
	DinoName   model.DinosaurName
	BaseHealth model.Health
	BaseMelee  model.Melee

	UniqueName       model.UniqueName
	HealthMultiplier model.UniqueMultiplier[model.Health]
	DamageMultiplier model.UniqueMultiplier[model.Melee]

	UniqueID model.UniqueDinosaurID
	Variants model.UniqueVariant
}

func NewCreateCreature(
	dinoName model.DinosaurName,
	baseHealth model.Health,
	baseMelee model.Melee,
	uniqueName model.UniqueName,
	healthMultiplier model.UniqueMultiplier[model.Health],
	damageMultiplier model.UniqueMultiplier[model.Melee],
	uniqueID model.UniqueDinosaurID,
	variants model.UniqueVariant,
) CreateCreature {
	return CreateCreature{
		DinoName:         dinoName,
		BaseHealth:       baseHealth,
		BaseMelee:        baseMelee,
		UniqueName:       uniqueName,
		HealthMultiplier: healthMultiplier,
		DamageMultiplier: damageMultiplier,
		UniqueID:         uniqueID,
		Variants:         variants,
	}
}

func (c CreateCreature) Dino() CreateDinosaur {
	return CreateDinosaur{
		name:       c.DinoName,
		baseHealth: c.BaseHealth,
		baseMelee:  c.BaseMelee,
	}
}

func (c CreateCreature) UniqueVariants() CreateVariants {
	return CreateVariants{
		variants: c.Variants,
	}
}

func (c CreateCreature) UniqueDinosaur(
	dinoID model.DinosaurID,
	uniqueVariantID model.UniqueVariantID,
) CreateUniqueDinosaur {
	return CreateUniqueDinosaur{
		name:             c.UniqueName,
		healthMultiplier: c.HealthMultiplier,
		damageMultiplier: c.DamageMultiplier,
		dinosaurID:       dinoID,
		uniqueVariantID:  uniqueVariantID,
	}
}

type CreateUniqueDinosaur struct {
	name             model.UniqueName
	healthMultiplier model.UniqueMultiplier[model.Health]
	damageMultiplier model.UniqueMultiplier[model.Melee]
	dinosaurID       model.DinosaurID
	uniqueVariantID  model.UniqueVariantID
}

func NewCreateUniqueDinosaur(
	name model.UniqueName,
	healthMultiplier model.UniqueMultiplier[model.Health],
	damageMultiplier model.UniqueMultiplier[model.Melee],
	dinosaurID model.DinosaurID,
	uniqueVariantID model.UniqueVariantID,
) CreateUniqueDinosaur {
	return CreateUniqueDinosaur{
		name:             name,
		healthMultiplier: healthMultiplier,
		damageMultiplier: damageMultiplier,
		dinosaurID:       dinosaurID,
		uniqueVariantID:  uniqueVariantID,
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

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
	dinoID           model.DinosaurID
	dinoName         model.DinosaurName
	baseHealth       model.Health
	baseMelee        model.Melee
	uniqueID         model.UniqueDinosaurID
	uniqueName       model.UniqueName
	healthMultiplier model.UniqueMultiplier[model.Health]
	damageMultiplier model.UniqueMultiplier[model.Melee]
	variantsID       model.UniqueVariantID
	variants         model.UniqueVariant
}

func NewUpdateCreature(
	dinoID model.DinosaurID,
	dinoName model.DinosaurName,
	baseHealth model.Health,
	baseMelee model.Melee,
	uniqueDinoID model.UniqueDinosaurID,
	uniqueName model.UniqueName,
	healthMultiplier model.UniqueMultiplier[model.Health],
	damageMultiplier model.UniqueMultiplier[model.Melee],
	variantsID model.UniqueVariantID,
	variants model.UniqueVariant,
) UpdateCreature {
	return UpdateCreature{
		dinoID:           dinoID,
		dinoName:         dinoName,
		baseHealth:       baseHealth,
		baseMelee:        baseMelee,
		uniqueID:         uniqueDinoID,
		uniqueName:       uniqueName,
		healthMultiplier: healthMultiplier,
		damageMultiplier: damageMultiplier,
		variantsID:       variantsID,
		variants:         variants,
	}
}

func (c UpdateCreature) Dino() UpdateDinosaur {
	return UpdateDinosaur{
		id:         c.dinoID,
		name:       c.dinoName,
		baseHealth: c.baseHealth,
		baseMelee:  c.baseMelee,
	}
}

func (c UpdateCreature) Unique() UpdateUniqueDinosaur {
	return UpdateUniqueDinosaur{
		uniqueDinoID:     c.uniqueID,
		name:             c.uniqueName,
		healthMultiplier: c.healthMultiplier,
		damageMultiplier: c.damageMultiplier,
	}
}

func (c UpdateCreature) Variants() UpdateVariants {
	return UpdateVariants{
		variantID: c.variantsID,
		variants:  c.variants,
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

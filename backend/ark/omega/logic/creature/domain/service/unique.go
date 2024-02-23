package service

import (
	"context"

	"github.com/samber/lo"

	"mods-explore/ark/omega/logic/creature/domain/model"
	variantModel "mods-explore/ark/omega/logic/variant/domain/model"
)

// UniqueQueryRepository 集約内のテーブルをjoinしてレコードを取得する処理を定義
type UniqueQueryRepository interface {
	Select(context.Context, model.UniqueDinosaurID) (*ResponseCreature, error)
	List(context.Context) (ResponseCreatures, error)
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
) CreateUniqueDinosaur {
	return CreateUniqueDinosaur{
		name:             c.UniqueName,
		healthMultiplier: c.HealthMultiplier,
		damageMultiplier: c.DamageMultiplier,
		dinosaurID:       dinoID,
	}
}

type CreateUniqueDinosaur struct {
	name             model.UniqueName
	healthMultiplier model.UniqueMultiplier[model.Health]
	damageMultiplier model.UniqueMultiplier[model.Melee]
	dinosaurID       model.DinosaurID
}

func NewCreateUniqueDinosaur(
	name model.UniqueName,
	healthMultiplier model.UniqueMultiplier[model.Health],
	damageMultiplier model.UniqueMultiplier[model.Melee],
	dinosaurID model.DinosaurID,
) CreateUniqueDinosaur {
	return CreateUniqueDinosaur{
		name:             name,
		healthMultiplier: healthMultiplier,
		damageMultiplier: damageMultiplier,
		dinosaurID:       dinosaurID,
	}
}

func (d CreateUniqueDinosaur) Name() model.UniqueName       { return d.name }
func (d CreateUniqueDinosaur) DinosaurID() model.DinosaurID { return d.dinosaurID }
func (d CreateUniqueDinosaur) HealthMultiplier() model.UniqueMultiplier[model.Health] {
	return d.healthMultiplier
}
func (d CreateUniqueDinosaur) DamageMultiplier() model.UniqueMultiplier[model.Melee] {
	return d.damageMultiplier
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
		variants: c.variants,
	}
}

type UpdateUniqueDinosaur struct {
	uniqueDinoID     model.UniqueDinosaurID
	dinosaurID       model.DinosaurID
	name             model.UniqueName
	healthMultiplier model.UniqueMultiplier[model.Health]
	damageMultiplier model.UniqueMultiplier[model.Melee]
}

func (d UpdateUniqueDinosaur) ID() model.UniqueDinosaurID   { return d.uniqueDinoID }
func (d UpdateUniqueDinosaur) Name() model.UniqueName       { return d.name }
func (d UpdateUniqueDinosaur) DinosaurID() model.DinosaurID { return d.dinosaurID }
func (d UpdateUniqueDinosaur) HealthMultiplier() model.UniqueMultiplier[model.Health] {
	return d.healthMultiplier
}
func (d UpdateUniqueDinosaur) DamageMultiplier() model.UniqueMultiplier[model.Melee] {
	return d.damageMultiplier
}

func NewUpdateUniqueDinosaur(
	uniqueDinoID model.UniqueDinosaurID,
	DinosaurID model.DinosaurID,
	name model.UniqueName,
	healthMultiplier model.UniqueMultiplier[model.Health],
	damageMultiplier model.UniqueMultiplier[model.Melee],
) UpdateUniqueDinosaur {
	return UpdateUniqueDinosaur{
		uniqueDinoID:     uniqueDinoID,
		dinosaurID:       DinosaurID,
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

type ResponseCreature struct {
	ResponseDinosaur
	ResponseVariants
	ResponseUnique
}

func (c ResponseCreature) ToUniqueDinosaur() model.UniqueDinosaur {
	vs := lo.Map(c.ResponseVariants.Values(), func(item model.DinosaurVariant, _ int) model.DinosaurVariant {
		return model.NewDinosaurVariant(
			variantModel.NewVariant(item.ID(), item.Group(), item.Name()),
			model.VariantDescriptions{},
		)
	})
	return model.NewUniqueDinosaur(
		model.NewDinosaur(
			c.ResponseDinosaur.ID(), c.ResponseDinosaur.Name(),
			c.ResponseDinosaur.Health(), c.ResponseDinosaur.Melee(),
		),
		c.ResponseUnique.ID(), c.ResponseUnique.Name(),
		c.ResponseUnique.HealthMultiplier(), c.ResponseUnique.MeleeMultiplier(),
		c.ResponseVariants.id, vs,
	)
}

type ResponseCreatures []ResponseCreature

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

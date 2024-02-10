package service

import (
	"context"

	"mods-explore/ark/omega/logic/variant/domain/model"
)

type CreateVariantGroup struct {
	name model.VariantGroupName
}

func NewCreateVariantGroup(name model.VariantGroupName) CreateVariantGroup {
	return CreateVariantGroup{name}
}
func (v CreateVariantGroup) Name() model.VariantGroupName { return v.name }

type UpdateVariantGroup struct {
	id   model.VariantGroupID
	name model.VariantGroupName
}

func NewUpdateVariantGroup(id model.VariantGroupID, name model.VariantGroupName) UpdateVariantGroup {
	return UpdateVariantGroup{id, name}
}

func (v UpdateVariantGroup) ID() model.VariantGroupID     { return v.id }
func (v UpdateVariantGroup) Name() model.VariantGroupName { return v.name }

type VariantGroupRepository interface {
	Select(context.Context, model.VariantID) (*model.Variant, error)
	List(context.Context) (model.Variants, error)
	Insert(context.Context, CreateVariant) (*model.Variant, error)
	Update(context.Context, UpdateVariantGroup) (*model.Variant, error)
	Delete(context.Context, model.VariantID) error
}

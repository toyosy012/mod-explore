package service

import (
	"context"

	"mods-explore/ark/omega/logic/variant/domain/model"
)

type CreateVariant struct {
	groupID model.VariantGroupID
	name    model.Name
}

func NewCreateVariant(groupID model.VariantGroupID, name model.Name) CreateVariant {
	return CreateVariant{groupID, name}
}
func (v CreateVariant) GroupID() model.VariantGroupID { return v.groupID }
func (v CreateVariant) Name() model.Name              { return v.name }

type UpdateVariant struct {
	id      model.VariantID
	groupID model.VariantGroupID
	name    model.Name
}

func NewUpdateVariant(id model.VariantID, groupID model.VariantGroupID, name model.Name) UpdateVariant {
	return UpdateVariant{id, groupID, name}
}

func (v UpdateVariant) ID() model.VariantID           { return v.id }
func (v UpdateVariant) GroupID() model.VariantGroupID { return v.groupID }
func (v UpdateVariant) Name() model.Name              { return v.name }

type VariantRepository interface {
	FindVariant(context.Context, model.VariantID) (*model.Variant, error)
	ListVariants(context.Context) (model.Variants, error)
	CreateVariant(context.Context, CreateVariant) (*model.Variant, error)
	UpdateVariant(context.Context, UpdateVariant) (*model.Variant, error)
	DeleteVariant(context.Context, model.VariantID) error
}

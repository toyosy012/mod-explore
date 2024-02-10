package service

import (
	"context"
	"errors"

	"mods-explore/ark/omega/logic/variant/domain/model"
)

var (
	NotFound            = errors.New("not found")
	IntervalServerError = errors.New("interval server error")
)

type CreateVariant struct {
	groupID model.GroupID
	name    model.Name
}

func NewCreateVariant(groupID model.GroupID, name model.Name) CreateVariant {
	return CreateVariant{groupID, name}
}
func (v CreateVariant) GroupID() model.GroupID { return v.groupID }
func (v CreateVariant) Name() model.Name       { return v.name }

type UpdateVariant struct {
	id      model.VariantID
	groupID model.GroupID
	name    model.Name
}

func NewUpdateVariant(id model.VariantID, groupID model.GroupID, name model.Name) UpdateVariant {
	return UpdateVariant{id, groupID, name}
}

func (v UpdateVariant) ID() model.VariantID    { return v.id }
func (v UpdateVariant) GroupID() model.GroupID { return v.groupID }
func (v UpdateVariant) Name() model.Name       { return v.name }

type VariantRepository interface {
	FindVariant(context.Context, model.VariantID) (*model.Variant, error)
	ListVariants(context.Context) (model.Variants, error)
	CreateVariant(context.Context, CreateVariant) (*model.Variant, error)
	UpdateVariant(context.Context, UpdateVariant) (*model.Variant, error)
	DeleteVariant(context.Context, model.VariantID) error
}

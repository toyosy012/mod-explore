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
	group model.GroupName
	name  model.Name
}

func NewCreateVariant(group model.GroupName, name model.Name) CreateVariant {
	return CreateVariant{group, name}
}
func (v CreateVariant) Group() model.GroupName { return v.group }
func (v CreateVariant) Name() model.Name       { return v.name }

type UpdateVariant struct {
	id    model.VariantID
	group model.GroupName
	name  model.Name
}

func NewUpdateVariant(id model.VariantID, group model.GroupName, name model.Name) UpdateVariant {
	return UpdateVariant{id, group, name}
}
func (v UpdateVariant) Group() model.GroupName { return v.group }
func (v UpdateVariant) Name() model.Name       { return v.name }

type VariantRepository interface {
	FindVariant(context.Context, model.VariantID) (*model.Variant, error)
	ListVariants(context.Context) (model.Variants, error)
	CreateVariant(context.Context, CreateVariant) (*model.Variant, error)
	UpdateVariant(context.Context, UpdateVariant) (*model.Variant, error)
	DeleteVariant(context.Context, model.VariantID) error
}

package service

import (
	"context"

	"mods-explore/ark/omega/logic/creature/domain/model"
)

type VariantsCommandRepository interface {
	Insert(context.Context, CreateVariants) (model.UniqueVariantID, error)
	Update(context.Context, UpdateVariants) error
	Delete(context.Context, model.UniqueVariantID) error
}

type CreateVariants struct {
	variants model.UniqueVariant
}

func NewCreateVariants(variants model.UniqueVariant) CreateVariants {
	return CreateVariants{
		variants: variants,
	}
}

type UpdateVariants struct {
	variantID model.UniqueVariantID
	variants  model.UniqueVariant
}

func NewUpdateVariants(
	variantID model.UniqueVariantID,
	variants model.UniqueVariant,
) UpdateVariants {
	return UpdateVariants{
		variantID: variantID,
		variants:  variants,
	}
}

type ResponseVariants struct {
	id       model.UniqueVariantID
	variants []model.DinosaurVariant
}

func NewResponseVariants(
	variantID model.UniqueVariantID,
	variants []model.DinosaurVariant,
) ResponseVariants {
	return ResponseVariants{id: variantID, variants: variants}
}

func (v ResponseVariants) ID() model.UniqueVariantID       { return v.id }
func (v ResponseVariants) Values() []model.DinosaurVariant { return v.variants }

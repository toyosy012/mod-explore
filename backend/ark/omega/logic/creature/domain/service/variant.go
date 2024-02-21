package service

import (
	"context"

	"github.com/samber/lo"
	variantModel "mods-explore/ark/omega/logic/variant/domain/model"

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

func (v CreateVariants) VariantIDs() []variantModel.VariantID {
	return lo.Map(v.variants, func(variant model.DinosaurVariant, _ int) variantModel.VariantID {
		return variant.ID()
	})
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

func (v UpdateVariants) VariantIDs() []variantModel.VariantID {
	return lo.Map(v.variants, func(variant model.DinosaurVariant, _ int) variantModel.VariantID {
		return variant.ID()
	})
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

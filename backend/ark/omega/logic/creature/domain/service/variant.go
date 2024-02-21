package service

import (
	"context"

	"github.com/samber/lo"
	variantModel "mods-explore/ark/omega/logic/variant/domain/model"

	"mods-explore/ark/omega/logic/creature/domain/model"
)

type UniqueVariantsCommand interface {
	Insert(context.Context, CreateVariants) (model.UniqueVariantID, error)
	Update(context.Context, UpdateVariants) error
	Delete(context.Context, model.UniqueVariantID) error
}

type CreateVariants struct {
	uniqueDinosaurID model.UniqueDinosaurID
	variants         model.UniqueVariant
}

func NewCreateVariants(id model.UniqueDinosaurID, variants model.UniqueVariant) CreateVariants {
	return CreateVariants{
		uniqueDinosaurID: id,
		variants:         variants,
	}
}

func (v CreateVariants) UniqueDinosaurID() model.UniqueDinosaurID { return v.uniqueDinosaurID }

func (v CreateVariants) VariantIDs() []variantModel.VariantID {
	return lo.Map(v.variants, func(variant model.DinosaurVariant, _ int) variantModel.VariantID {
		return variant.ID()
	})
}

type UpdateVariants struct {
	uniqueDinosaurID model.UniqueDinosaurID
	variants         model.UniqueVariant
}

func NewUpdateVariants(
	uniqueDinosaurID model.UniqueDinosaurID,
	variants model.UniqueVariant,
) UpdateVariants {
	return UpdateVariants{
		uniqueDinosaurID: uniqueDinosaurID,
		variants:         variants,
	}
}

func (v UpdateVariants) UniqueDinosaurID() model.UniqueDinosaurID { return v.uniqueDinosaurID }
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

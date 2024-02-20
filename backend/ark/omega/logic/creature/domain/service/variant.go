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
	uniqueID model.UniqueDinosaurID
	variants model.UniqueVariant
}

func NewCreateVariants(uniqueID model.UniqueDinosaurID, variants model.UniqueVariant) CreateVariants {
	return CreateVariants{
		uniqueID: uniqueID,
		variants: variants,
	}
}

type UpdateVariants struct {
	variantID model.UniqueVariantID
	uniqueID  model.UniqueDinosaurID
	variants  model.UniqueVariant
}

func NewUpdateVariants(
	variantID model.UniqueVariantID,
	uniqueID model.UniqueDinosaurID,
	variants model.UniqueVariant,
) UpdateVariants {
	return UpdateVariants{
		variantID: variantID,
		uniqueID:  uniqueID,
		variants:  variants,
	}
}

type ResponseVariants struct {
	variantID model.UniqueVariantID
	uniqueID  model.UniqueDinosaurID
	variants  []model.DinosaurVariant
}

func NewResponseVariants(
	variantID model.UniqueVariantID,
	uniqueID model.UniqueDinosaurID,
	variants []model.DinosaurVariant,
) ResponseVariants {
	return ResponseVariants{variantID: variantID, uniqueID: uniqueID, variants: variants}
}

func (v ResponseVariants) ID() model.UniqueVariantID        { return v.variantID }
func (v ResponseVariants) UniqueID() model.UniqueDinosaurID { return v.uniqueID }
func (v ResponseVariants) Values() []model.DinosaurVariant  { return v.variants }

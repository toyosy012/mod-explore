package service

import (
	"context"

	variantModel "mods-explore/ark/omega/logic/variant/domain/model"

	"mods-explore/ark/omega/logic/creature/domain/model"
)

type UniqueVariantsCommand interface {
	Insert(context.Context, CreateVariants) error
	Update(context.Context, UpdateVariants) error
	Delete(context.Context, model.UniqueVariantID) error
}

type CreateVariants struct {
	uniqueDinosaurID model.UniqueDinosaurID
	variantIDs       [2]variantModel.VariantID
}

func (v CreateVariants) UniqueDinosaurID() model.UniqueDinosaurID { return v.uniqueDinosaurID }

func (v CreateVariants) VariantIDs() [2]variantModel.VariantID {
	return v.variantIDs
}

type UpdateVariants struct {
	uniqueDinosaurID model.UniqueDinosaurID
	variantIDs       [2]variantModel.VariantID
}

func (v UpdateVariants) UniqueDinosaurID() model.UniqueDinosaurID { return v.uniqueDinosaurID }
func (v UpdateVariants) VariantIDs() [2]variantModel.VariantID {
	return v.variantIDs
}

type ResponseVariants struct {
	variants [2]model.DinosaurVariant
}

func NewResponseVariants(
	variants [2]model.DinosaurVariant,
) ResponseVariants {
	return ResponseVariants{variants: variants}
}

func (v ResponseVariants) Values() [2]model.DinosaurVariant { return v.variants }

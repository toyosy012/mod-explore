package service

import (
	"context"

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
	id       model.UniqueVariantID
	variants [2]model.DinosaurVariant
}

func NewResponseVariants(
	variantID model.UniqueVariantID,
	variants [2]model.DinosaurVariant,
) ResponseVariants {
	return ResponseVariants{id: variantID, variants: variants}
}

func (v ResponseVariants) ID() model.UniqueVariantID        { return v.id }
func (v ResponseVariants) Values() [2]model.DinosaurVariant { return v.variants }

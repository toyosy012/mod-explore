package model

import (
	"mods-explore/ark/omega/logic/variant/domain/model"
)

type UniqueVariantID int
type VariantDescription string
type VariantDescriptions []VariantDescription
type DinosaurVariant struct {
	model.Variant
	descriptions VariantDescriptions
}

func NewDinosaurVariant(
	variant model.Variant,
	descriptions VariantDescriptions,
) DinosaurVariant {
	return DinosaurVariant{
		Variant:      variant,
		descriptions: descriptions,
	}
}

func (v DinosaurVariant) Descriptions() VariantDescriptions { return v.descriptions }

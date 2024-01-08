package model

import (
	"errors"

	"mods-explore/ark/omega/logic/variant/domain/model"
)

var (
	minMultiplier    float32 = 0.0
	errMinMultiplier         = errors.New("バリアントグループの倍率補正の値を0より大きくしてください")
)

type VariantGroupMultiplier float32

func (m VariantGroupMultiplier) ToFloat32() float32 { return float32(m) }

func NewVariantGroupMultiplier(value float32) (VariantGroupMultiplier, error) {
	if minMultiplier >= value {
		return VariantGroupMultiplier(1.0), errMinMultiplier
	}
	return VariantGroupMultiplier(value), nil
}

type VariantDescription string
type VariantDescriptions []VariantDescription
type DinosaurVariant struct {
	model.Variant
	multiplier   VariantGroupMultiplier
	descriptions VariantDescriptions
}

func NewDinosaurVariant(
	variant model.Variant,
	multiplier VariantGroupMultiplier,
	descriptions VariantDescriptions,
) DinosaurVariant {
	return DinosaurVariant{
		Variant:      variant,
		multiplier:   multiplier,
		descriptions: descriptions,
	}
}

func (v DinosaurVariant) GroupMultiplier() VariantGroupMultiplier { return v.multiplier }
func (v DinosaurVariant) Descriptions() VariantDescriptions       { return v.descriptions }

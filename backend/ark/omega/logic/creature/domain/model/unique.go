package model

import (
	"errors"
)

var (
	errUniqueMinMultiplier float32 = 0.0
)

type UniqueDinosaur struct {
	Dinosaur

	uniqueDinoID     UniqueDinosaurID
	uniqueName       UniqueName
	variants         UniqueVariant
	healthMultiplier UniqueMultiplier[Health]
	damageMultiplier UniqueMultiplier[Melee]
}

func NewUniqueDinosaur(
	base Dinosaur,
	id UniqueDinosaurID,
	name UniqueName,
	variants UniqueVariant,
	healthMultiplier UniqueMultiplier[Health],
	damageMultiplier UniqueMultiplier[Melee],
) UniqueDinosaur {
	return UniqueDinosaur{
		Dinosaur:         base,
		uniqueDinoID:     id,
		uniqueName:       name,
		variants:         variants,
		healthMultiplier: healthMultiplier,
		damageMultiplier: damageMultiplier,
	}
}

type UniqueDinosaurID int
type UniqueName string

// DinosaurStatus multiplierでfloat32との計算に用いるため、数値型のみに限定する
type DinosaurStatus interface {
	Health | Melee
}

type UniqueMultiplier[T DinosaurStatus] struct{ value UniqueTotalMultiplier }

func NewUniqueMultiplier[T DinosaurStatus](v UniqueTotalMultiplier) (*UniqueMultiplier[T], error) {
	if errUniqueMinMultiplier >= v.ToFloat32() {
		return nil, errors.New("ユニーク生物のステータス倍率は0より大きくしてください")
	}
	return &UniqueMultiplier[T]{value: v}, nil
}

type UniqueMultipliedStatus[T DinosaurStatus] float32

// multiple UniqueMultiplierに与えた型引数と同じ型のbaseを与えないとエラーになるようにする
func (u UniqueMultiplier[T]) multiple(base T) UniqueMultipliedStatus[T] {
	return UniqueMultipliedStatus[T](float32(base) * u.value.ToFloat32())
}

func (d UniqueDinosaur) Health() UniqueMultipliedStatus[Health] {
	return d.healthMultiplier.multiple(d.baseHealth)
}

func (d UniqueDinosaur) Damage() UniqueMultipliedStatus[Melee] {
	return d.damageMultiplier.multiple(d.baseMelee)
}

type UniqueVariant [2]DinosaurVariant

type UniqueTotalMultiplier float32

func (m UniqueTotalMultiplier) ToFloat32() float32 { return float32(m) }

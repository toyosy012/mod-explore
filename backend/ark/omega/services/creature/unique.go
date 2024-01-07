package creature

import (
	"errors"

	"github.com/samber/lo"
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
		return nil, errors.New("ユニークバリアントの合計倍率は0より大きくしてください")
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

func (v UniqueVariant) TotalMultiplier() UniqueTotalMultiplier {
	return lo.ReduceRight[DinosaurVariant, UniqueTotalMultiplier](
		v[:], // slice to list
		func(agg UniqueTotalMultiplier, item DinosaurVariant, _ int) UniqueTotalMultiplier {
			return UniqueTotalMultiplier(float32(agg) * item.GroupMultiplier().ToFloat32())
		},
		1.0,
	)
}

type UniqueTotalMultiplier float32

func (m UniqueTotalMultiplier) ToFloat32() float32 { return float32(m) }

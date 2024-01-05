package creature

import "errors"

type UniqueDinosaur struct {
	Dinosaur

	uniqueDinoID     UniqueDinosaurID
	uniqueName       UniqueName
	variants         UniqueVariants
	healthMultiplier UniqueMultiplier[Health]
	damageMultiplier UniqueMultiplier[Melee]
}

func NewUniqueDinosaur(
	base Dinosaur,
	id UniqueDinosaurID,
	name UniqueName,
	variants UniqueVariants,
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

type UniqueVariants [2]string

type UniqueMultipliedStatus interface {
	Health | Melee
}

type UniqueMultiplier[T UniqueMultipliedStatus] struct{ value T }

func NewUniqueMultiplier[T UniqueMultipliedStatus](v T) (*UniqueMultiplier[T], error) {
	if 0 == v {
		return nil, errors.New("倍率は0より大きくする必要があります")
	}
	return &UniqueMultiplier[T]{value: v}, nil
}

// multiple UniqueMultiplierに与えた型引数と同じ型のbaseを与えないとエラーになるようにする
func (u UniqueMultiplier[T]) multiple(base T) T { return base * u.value }

func (d UniqueDinosaur) Health() Health {
	return d.healthMultiplier.multiple(d.baseHealth)
}

func (d UniqueDinosaur) Damage() Melee {
	return d.damageMultiplier.multiple(d.baseMelee)
}

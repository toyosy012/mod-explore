package model

import (
	"errors"
)

type Health uint

func NewHealth(value uint) (Health, error) {
	if 0 == value {
		return 0, errors.New("体力0は許容されない不正な値です")
	}
	return Health(value), nil
}

type Melee uint

// NewMelee 攻撃不可な生物も存在するのでその場合は0を指定
func NewMelee(value uint) Melee { return Melee(value) }

type DinosaurID int
type DinosaurName string

type Dinosaur struct {
	id         DinosaurID
	name       DinosaurName
	baseHealth Health
	baseMelee  Melee
}

func NewDinosaur(id DinosaurID, name DinosaurName, health Health, Melee Melee) Dinosaur {
	return Dinosaur{
		id:         id,
		name:       name,
		baseHealth: health,
		baseMelee:  Melee,
	}
}

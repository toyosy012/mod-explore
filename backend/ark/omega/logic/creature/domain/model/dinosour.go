package model

import (
	"errors"
)

type Health uint

func (h Health) Value() uint { return uint(h) }

func NewHealth(value uint) (Health, error) {
	if 0 == value {
		return 0, errors.New("体力0は許容されない不正な値です")
	}
	return Health(value), nil
}

type Melee uint

func (m Melee) Value() uint { return uint(m) }

// NewMelee 攻撃不可な生物も存在するのでその場合は0を指定
func NewMelee(value uint) Melee { return Melee(value) }

type DinosaurID int

func (i DinosaurID) Value() int { return int(i) }

type DinosaurName string

func (n DinosaurName) Value() string { return string(n) }

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

func (d Dinosaur) BaseID() DinosaurID     { return d.id }
func (d Dinosaur) BaseName() DinosaurName { return d.name }
func (d Dinosaur) Health() Health         { return d.baseHealth }
func (d Dinosaur) Melee() Melee           { return d.baseMelee }

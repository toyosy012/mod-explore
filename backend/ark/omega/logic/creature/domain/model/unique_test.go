package model

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"mods-explore/ark/omega/logic/variant/domain/model"
)

type UniqueDinosaurTestSuite struct {
	suite.Suite

	baseDino                Dinosaur
	defaultID               UniqueDinosaurID
	defaultName             UniqueName
	defaultVariantsID       UniqueVariantID
	variants                UniqueVariant
	defaultHealthMultiplier UniqueMultiplier[Health]
	defaultDamageMultiplier UniqueMultiplier[Melee]
}

const (
	minHealth              = 0
	healthValue            = 2
	healthUniqueMultiplier = 36.0
	multipliedHealth       = 72.0
	minMelee               = 0
	meleeValue             = 2
	meleeUniqueMultiplier  = 36.0
	multipliedMelee        = 72.0
)

func TestUniqueDinosaur(t *testing.T) {
	s, err := NewUniqueDinosaurTestSuite()
	if err != nil {
		t.Fatal(err)
	}
	suite.Run(t, s)
}

func NewUniqueDinosaurTestSuite() (*UniqueDinosaurTestSuite, error) {
	baseHealth, err := NewHealth(healthValue)
	if err != nil {
		return nil, err
	}
	baseMelee := NewMelee(meleeValue)

	dino := NewDinosaur(
		DinosaurID(1),
		"Dodo",
		baseHealth,
		baseMelee,
	)

	variants := UniqueVariant(
		[]DinosaurVariant{
			NewDinosaurVariant(
				model.NewVariant(cosmicID, cosmic, singularity),
				[]VariantDescription{
					"AoE explosive tick damage, traps dinos in center.",
					"Destroys corpses.",
				},
			),
			NewDinosaurVariant(
				model.NewVariant(natureID, nature, thunderstorm),
				[]VariantDescription{
					"Summons lightning bolts within an area to strike random targets.",
				},
			),
		},
	)

	defaultHealthMultiplier, err := NewUniqueMultiplier[Health](healthUniqueMultiplier)
	if err != nil {
		return nil, err
	}
	defaultDamageMultiplier, err := NewUniqueMultiplier[Melee](meleeUniqueMultiplier)
	if err != nil {
		return nil, err
	}

	return &UniqueDinosaurTestSuite{
		baseDino: dino,

		defaultID:               UniqueDinosaurID(1),
		defaultName:             "Kenny",
		defaultVariantsID:       UniqueVariantID(1),
		variants:                variants,
		defaultHealthMultiplier: *defaultHealthMultiplier,
		defaultDamageMultiplier: *defaultDamageMultiplier,
	}, nil
}

func (s *UniqueDinosaurTestSuite) TestMultiplierHealth() {
	s.T().Log("体力型で倍率の型とベース値の計算が可能かテスト")

	uniqueDino := NewUniqueDinosaur(
		s.baseDino, UniqueDinosaurID(1), s.defaultName,
		s.defaultHealthMultiplier, s.defaultDamageMultiplier,
		s.defaultVariantsID, s.variants,
	)

	uniqueHealth := uniqueDino.Health()
	health := UniqueMultipliedStatus[Health](multipliedHealth)
	s.Equal(health, uniqueHealth)
}

func (s *UniqueDinosaurTestSuite) TestErrMultiplierHealthZero() {
	s.T().Log("体力型で倍率が0のエラーケースのテスト")

	if _, err := NewUniqueMultiplier[Health](minHealth); err == nil {
		s.T().Errorf("体力型の倍率が0でエラーになっていません")
	}
}

func (s *UniqueDinosaurTestSuite) TestMultiplierDamage() {
	s.T().Log("攻撃力型で倍率の型とベース値の計算が可能かテスト")

	uniqueDino := NewUniqueDinosaur(
		s.baseDino, UniqueDinosaurID(1), "Kenny",
		s.defaultHealthMultiplier, s.defaultDamageMultiplier,
		s.defaultVariantsID, s.variants,
	)

	uniqueHealth := uniqueDino.Damage()
	melee := UniqueMultipliedStatus[Melee](multipliedMelee)
	s.Equal(melee, uniqueHealth)
}

func (s *UniqueDinosaurTestSuite) TestErrMultiplierDamageZero() {
	s.T().Log("攻撃力型で倍率が0のエラーケースのテス")

	if _, err := NewUniqueMultiplier[Melee](minMelee); err == nil {
		s.T().Errorf("攻撃力型の倍率が0でエラーになっていません")
	}
}

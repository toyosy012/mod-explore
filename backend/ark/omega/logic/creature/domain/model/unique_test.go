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
	variants                UniqueVariant
	defaultHealthMultiplier UniqueMultiplier[Health]
	defaultDamageMultiplier UniqueMultiplier[Melee]
}

func TestUniqueDinosaur(t *testing.T) {
	s, err := NewUniqueDinosaurTestSuite()
	if err != nil {
		t.Fatal(err)
	}
	suite.Run(t, s)
}

func NewUniqueDinosaurTestSuite() (*UniqueDinosaurTestSuite, error) {
	baseHealth, err := NewHealth(2)
	if err != nil {
		return nil, err
	}
	baseMelee := NewMelee(2)

	dino := NewDinosaur(
		DinosaurID(1),
		"Dodo",
		baseHealth,
		baseMelee,
	)

	cosmicMultiplier, err := NewVariantGroupMultiplier(6.0)
	if err != nil {
		return nil, err
	}
	natureMultiplier, err := NewVariantGroupMultiplier(6.0)
	if err != nil {
		return nil, err
	}

	variants := UniqueVariant(
		[2]DinosaurVariant{
			NewDinosaurVariant(
				model.NewVariant(cosmicID, cosmic, singularity),
				cosmicMultiplier,
				[]VariantDescription{
					"AoE explosive tick damage, traps dinos in center.",
					"Destroys corpses.",
				},
			),
			NewDinosaurVariant(
				model.NewVariant(natureID, nature, thunderstorm),
				natureMultiplier,
				[]VariantDescription{
					"Summons lightning bolts within an area to strike random targets.",
				},
			),
		},
	)

	defaultHealthMultiplier, err := NewUniqueMultiplier[Health](variants.TotalMultiplier())
	if err != nil {
		return nil, err
	}
	defaultDamageMultiplier, err := NewUniqueMultiplier[Melee](variants.TotalMultiplier())
	if err != nil {
		return nil, err
	}

	return &UniqueDinosaurTestSuite{
		baseDino: dino,

		defaultID:               UniqueDinosaurID(1),
		defaultName:             "Kenny",
		variants:                variants,
		defaultHealthMultiplier: *defaultHealthMultiplier,
		defaultDamageMultiplier: *defaultDamageMultiplier,
	}, nil
}

func (s *UniqueDinosaurTestSuite) TestTotalMultiplier() {
	s.T().Log("ユニーク生物の持つバリアント倍率が乗算で算出されているか")
	expect := UniqueTotalMultiplier(36.0)
	s.Equal(expect, s.variants.TotalMultiplier())
}

func (s *UniqueDinosaurTestSuite) TestMultiplierHealth() {
	s.T().Log("体力型で倍率の型とベース値の計算が可能かテスト")

	uniqueDino := NewUniqueDinosaur(
		s.baseDino, UniqueDinosaurID(1), s.defaultName, s.variants, s.defaultHealthMultiplier, s.defaultDamageMultiplier,
	)

	uniqueHealth := uniqueDino.Health()
	health := UniqueMultipliedStatus[Health](72.0)
	s.Equal(health, uniqueHealth)
}

func (s *UniqueDinosaurTestSuite) TestErrMultiplierHealthZero() {
	s.T().Log("体力型で倍率が0のエラーケースのテスト")

	if _, err := NewUniqueMultiplier[Health](0); err == nil {
		s.T().Errorf("体力型の倍率が0でエラーになっていません")
	}
}

func (s *UniqueDinosaurTestSuite) TestMultiplierDamage() {
	s.T().Log("攻撃力型で倍率の型とベース値の計算が可能かテスト")

	uniqueDino := NewUniqueDinosaur(
		s.baseDino, UniqueDinosaurID(1), "Kenny", s.variants, s.defaultHealthMultiplier, s.defaultDamageMultiplier,
	)

	uniqueHealth := uniqueDino.Damage()
	melee := UniqueMultipliedStatus[Melee](72.0)
	s.Equal(melee, uniqueHealth)
}

func (s *UniqueDinosaurTestSuite) TestErrMultiplierDamageZero() {
	s.T().Log("攻撃力型で倍率が0のエラーケースのテス")

	if _, err := NewUniqueMultiplier[Melee](0); err == nil {
		s.T().Errorf("攻撃力型の倍率が0でエラーになっていません")
	}
}
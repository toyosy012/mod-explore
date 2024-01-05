package creature

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UniqueDinosaurTest struct {
	suite.Suite

	baseDino                Dinosaur
	defaultID               UniqueDinosaurID
	defaultName             UniqueName
	variants                [2]string
	defaultHealthMultiplier UniqueMultiplier[Health]
	defaultDamageMultiplier UniqueMultiplier[Melee]
}

func TestUserAccountService(t *testing.T) {
	s := NewUniqueDinosaurTestSuite()
	suite.Run(t, &s)
}

func NewUniqueDinosaurTestSuite() UniqueDinosaurTest {
	baseHealth, _ := NewHealth(2)
	baseMelee := NewMelee(2)

	dino := NewDinosaur(
		DinosaurID(1),
		"Dodo",
		baseHealth,
		baseMelee,
	)

	defaultHealthMultiplier, _ := NewUniqueMultiplier[Health](1)
	defaultDamageMultiplier, _ := NewUniqueMultiplier[Melee](1)

	return UniqueDinosaurTest{
		baseDino: dino,

		defaultID:               UniqueDinosaurID(1),
		defaultName:             "Kenny",
		variants:                [2]string{"Phoenix", "Self-Destructive"},
		defaultHealthMultiplier: *defaultHealthMultiplier,
		defaultDamageMultiplier: *defaultDamageMultiplier,
	}
}

func (s *UniqueDinosaurTest) TestMultiplierHealth() {
	s.T().Log("体力型で倍率の型とベース値の計算が可能かテスト")

	healthMultiplier, _ := NewUniqueMultiplier[Health](5)
	uniqueDino := NewUniqueDinosaur(
		s.baseDino, UniqueDinosaurID(1), s.defaultName, s.variants, *healthMultiplier, s.defaultDamageMultiplier,
	)

	uniqueHealth := uniqueDino.healthMultiplier.multiple(uniqueDino.Dinosaur.baseHealth)
	health, _ := NewHealth(10)
	s.Equal(health, uniqueHealth)
}

func (s *UniqueDinosaurTest) TestErrMultiplierHealthZero() {
	s.T().Log("体力型で倍率が0のエラーケースのテスト")

	if _, err := NewUniqueMultiplier[Health](0); err == nil {
		s.T().Log("体力型の倍率が0でエラーになっていません")
	}
}

func (s *UniqueDinosaurTest) TestMultiplierDamage() {
	s.T().Log("攻撃力型で倍率の型とベース値の計算が可能かテスト")

	healthMultiplier, _ := NewUniqueMultiplier[Health](5)
	uniqueDino := NewUniqueDinosaur(
		s.baseDino, UniqueDinosaurID(1), "Kenny", s.variants, *healthMultiplier, s.defaultDamageMultiplier,
	)

	uniqueHealth := uniqueDino.healthMultiplier.multiple(uniqueDino.Dinosaur.baseHealth)
	health, _ := NewHealth(10)
	s.Equal(health, uniqueHealth)
}

func (s *UniqueDinosaurTest) TestErrMultiplierDamageZero() {
	s.T().Log("攻撃力型で倍率が0のエラーケースのテス")

	if _, err := NewUniqueMultiplier[Melee](0); err == nil {
		s.T().Log("攻撃力型の倍率が0でエラーになっていません")
	}
}

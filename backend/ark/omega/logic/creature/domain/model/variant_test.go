package model

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"mods-explore/ark/omega/logic/variant/domain/model"
)

var (
	cosmic   = model.VariantGroupName("Cosmic")
	cosmicID = model.VariantID(1)
	nature   = model.VariantGroupName("Nature")
	natureID = model.VariantID(2)
)

var (
	singularity  = model.Name("Singularity")
	thunderstorm = model.Name("Thunderstorm")
)

type DinosaurVariantTestSuite struct {
	suite.Suite

	dinoSingularity  DinosaurVariant
	dinoThunderstorm DinosaurVariant
}

func NewDinosaurVariantTestSuite() (*DinosaurVariantTestSuite, error) {
	return &DinosaurVariantTestSuite{
		dinoSingularity: NewDinosaurVariant(
			model.NewVariant(cosmicID, cosmic, singularity),
			[]VariantDescription{
				"AoE explosive tick damage, traps dinos in center.",
				"Destroys corpses.",
			},
		),
		dinoThunderstorm: NewDinosaurVariant(
			model.NewVariant(natureID, nature, thunderstorm),
			[]VariantDescription{
				"Summons lightning bolts within an area to strike random targets.",
			},
		),
	}, nil
}

func TestDinosaurVariant(t *testing.T) {
	s, err := NewDinosaurVariantTestSuite()
	if err != nil {
		t.Fatal(err)
	}
	suite.Run(t, s)
}

func (s *DinosaurVariantTestSuite) TestGetVariantGroupName() {
	s.T().Log("バリアントのグループ名を取得できるか")
	s.Equal(cosmic, s.dinoSingularity.Group())
}

func (s *DinosaurVariantTestSuite) TestGetVariantName() {
	s.T().Log("バリアント名を取得できるか")
	s.Equal(singularity, s.dinoSingularity.Name())
}

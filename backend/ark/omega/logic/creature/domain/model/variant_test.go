package model

import (
	"mods-explore/ark/omega/logic/variant/domain/model"
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	cosmic   = model.GroupName("Cosmic")
	cosmicID = model.VariantID(1)
	nature   = model.GroupName("Nature")
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
	cosmicMultiplier, err := NewVariantGroupMultiplier(6.0)
	if err != nil {
		return nil, err
	}
	natureMultiplier, err := NewVariantGroupMultiplier(6.0)
	if err != nil {
		return nil, err
	}
	return &DinosaurVariantTestSuite{
		dinoSingularity: NewDinosaurVariant(
			model.NewVariant(cosmicID, cosmic, singularity),
			cosmicMultiplier,
			[]VariantDescription{
				"AoE explosive tick damage, traps dinos in center.",
				"Destroys corpses.",
			},
		),
		dinoThunderstorm: NewDinosaurVariant(
			model.NewVariant(natureID, nature, thunderstorm),
			natureMultiplier,
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

func (s *DinosaurVariantTestSuite) TestInitGroupMultiplier() {
	s.T().Log("バリアントのグループ倍率の初期化処理のチェック")

	multiplier, err := NewVariantGroupMultiplier(1.0)
	if err != nil {
		s.T().Error("正常に生成されるべき処理で異常が発生しました。")
	}
	s.Equal(VariantGroupMultiplier(1.0), multiplier)
}

func (s *DinosaurVariantTestSuite) TestInitZeroGroupMultiplier() {
	s.T().Log("バリアントのグループ倍率が0の初期化処理のチェック")

	_, err := NewVariantGroupMultiplier(0.0)
	s.ErrorIsf(errMinMultiplier, err, "倍率0.0倍の初期化処理のエラー判定に失敗")
}

func (s *DinosaurVariantTestSuite) TestInitMinusGroupMultiplier() {
	s.T().Log("バリアントのグループ倍率が負の時の初期化処理のチェック")

	_, err := NewVariantGroupMultiplier(-0.1)
	s.ErrorIsf(errMinMultiplier, err, "負の倍率の初期化処理のエラー判定に失敗")
}

func (s *DinosaurVariantTestSuite) TestGetVariantGroupName() {
	s.T().Log("バリアントのグループ名を取得できるか")
	s.Equal(cosmic, s.dinoSingularity.Group())
}

func (s *DinosaurVariantTestSuite) TestGetVariantName() {
	s.T().Log("バリアント名を取得できるか")
	s.Equal(singularity, s.dinoSingularity.Name())
}

func (s *DinosaurVariantTestSuite) TestGetGroupMultiplier() {
	s.T().Log("バリアントの倍率を取得できるか")
	base := 6.0
	s.Equal(VariantGroupMultiplier(base), s.dinoSingularity.GroupMultiplier())
}

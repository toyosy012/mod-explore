package usecase

import (
	"errors"
	"testing"

	"github.com/morikuni/failure"
	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"mods-explore/ark/omega/logic"
	"mods-explore/ark/omega/logic/creature/domain/model"
	"mods-explore/ark/omega/logic/creature/domain/service"
	variantModel "mods-explore/ark/omega/logic/variant/domain/model"
)

const (
	successUniqueID = iota
	notExistUniqueID
	internalServerErrUniqueID
	errUniqueID
)

type UniqueDinosaurTestSuite struct {
	suite.Suite

	mockDB  *mockUniqueDB
	usecase UniqueUsecase
	unique  model.UniqueDinosaur
}

func newTestUniqueDinosaurSuite() *UniqueDinosaurTestSuite { return &UniqueDinosaurTestSuite{} }

func TestUniqueDinosaurSuite(t *testing.T) {
	suite.Run(t, newTestUniqueDinosaurSuite())
}

const (
	findUnique = "Select"
)

func (s *UniqueDinosaurTestSuite) SetupSuite() {
	injector := do.New()

	mockDB := newMockUniqueDB()
	do.ProvideValue[service.UniqueRepository](injector, mockDB)
	s.mockDB = mockDB
	usecase, err := NewUnique(injector)
	if err != nil {
		return
	}

	s.usecase = usecase
}

const (
	multiplierHealth = 36.0
	health           = 1
	multiplierMelee  = 36.0
	melee            = 0
	cosmicID         = iota
	natureID
	creatureID = 1
	uniqueID   = 1
)

func (s *UniqueDinosaurTestSuite) BeforeTest(_ string, tableName string) {
	if "TestFind" == tableName {
		h, err := model.NewHealth(health)
		if err != nil {
			s.T().Error(err)
			return
		}
		healthMultiplier, err := model.NewUniqueMultiplier[model.Health](multiplierHealth)
		if err != nil {
			s.T().Error(err)
			return
		}
		meleeMultiplier, err := model.NewUniqueMultiplier[model.Melee](multiplierMelee)
		if err != nil {
			s.T().Error(err)
			return
		}
		s.unique = model.NewUniqueDinosaur(
			model.NewDinosaur(creatureID, "dodo", h, model.NewMelee(melee)),
			uniqueID,
			"Kenny",
			model.UniqueVariant{
				model.NewDinosaurVariant(
					variantModel.NewVariant(cosmicID, "cosmic", "singularity"),
					model.VariantDescriptions{},
				),
				model.NewDinosaurVariant(
					variantModel.NewVariant(natureID, "nature", "thunder"),
					model.VariantDescriptions{},
				),
			},
			*healthMultiplier,
			*meleeMultiplier,
		)
	}
}

func (s *UniqueDinosaurTestSuite) TestFind() {
	{
		s.mockDB.On(
			findVariantGroup,
			ctx,
			model.UniqueDinosaurID(successUniqueID),
		).
			Return(&s.unique, nil)
		r, err := s.usecase.Find(ctx, model.UniqueDinosaurID(successUniqueID))
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(&s.unique, r)
	}
	{
		s.mockDB.On(
			findUnique,
			ctx,
			model.UniqueDinosaurID(notExistUniqueID),
		).
			Return(nil, service.NotFound)
		_, err := s.usecase.Find(ctx, model.UniqueDinosaurID(notExistUniqueID))
		s.True(failure.Is(err, logic.NotFound))
	}
	{
		s.mockDB.On(
			findUnique,
			ctx,
			model.UniqueDinosaurID(internalServerErrUniqueID),
		).
			Return(nil, service.IntervalServerError)
		_, err := s.usecase.Find(ctx, model.UniqueDinosaurID(internalServerErrUniqueID))
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDB.On(
			findUnique,
			ctx,
			model.UniqueDinosaurID(errUniqueID),
		).
			Return(nil, failure.Wrap(e))
		_, err := s.usecase.Find(ctx, model.UniqueDinosaurID(errUniqueID))
		s.True(errors.Is(err, e))
	}
}

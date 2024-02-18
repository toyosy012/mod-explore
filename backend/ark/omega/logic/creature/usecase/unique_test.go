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
	create  service.CreateUniqueDinosaur
	update  service.UpdateUniqueDinosaur
}

func newTestUniqueDinosaurSuite() *UniqueDinosaurTestSuite { return &UniqueDinosaurTestSuite{} }

func TestUniqueDinosaurSuite(t *testing.T) {
	suite.Run(t, newTestUniqueDinosaurSuite())
}

const (
	findUnique   = "Select"
	listUnique   = "List"
	insertUnique = "Insert"
	updateUnique = "Update"
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
	variants := model.UniqueVariant{
		model.NewDinosaurVariant(
			variantModel.NewVariant(cosmicID, cosmic, singularity),
			model.VariantDescriptions{},
		),
		model.NewDinosaurVariant(
			variantModel.NewVariant(natureID, nature, thunderstorm),
			model.VariantDescriptions{},
		),
	}
	{
		s.unique = model.NewUniqueDinosaur(
			model.NewDinosaur(creatureID, creatureName, h, model.NewMelee(melee)),
			uniqueID, uniqueName, variants, *healthMultiplier, *meleeMultiplier,
		)
	}
	{
		s.create = service.NewCreateUniqueDinosaur(
			service.NewCreateDinosaur(creatureName, h, melee),
			uniqueName, variants, *healthMultiplier, *meleeMultiplier,
		)
	}
	{
		s.update = service.NewUpdateUniqueDinosaur(
			service.NewUpdateDinosaur(creatureID, creatureName, h, melee),
			uniqueID, uniqueName, variants, *healthMultiplier, *meleeMultiplier,
		)
	}
}

const (
	multiplierHealth = 36.0
	health           = 1
	multiplierMelee  = 36.0
	melee            = 0
	cosmicID         = iota
	natureID
	creatureID   = 1
	creatureName = "dodo"
	uniqueID     = 1
	uniqueName   = "Kenny"
	cosmic       = "cosmic"
	singularity  = "singularity"
	nature       = "nature"
	thunderstorm = "thunderstorm"
)

func (s *UniqueDinosaurTestSuite) TestFind() {
	{
		s.mockDB.On(
			findUnique,
			ctx,
			model.UniqueDinosaurID(successUniqueID),
		).
			Return(&s.unique, nil).
			Once()
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
			Return(nil, service.NotFound).
			Once()
		_, err := s.usecase.Find(ctx, model.UniqueDinosaurID(notExistUniqueID))
		s.True(failure.Is(err, logic.NotFound))
	}
	{
		s.mockDB.On(
			findUnique,
			ctx,
			model.UniqueDinosaurID(internalServerErrUniqueID),
		).
			Return(nil, service.IntervalServerError).
			Once()
		_, err := s.usecase.Find(ctx, model.UniqueDinosaurID(internalServerErrUniqueID))
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDB.On(
			findUnique,
			ctx,
			model.UniqueDinosaurID(errUniqueID),
		).
			Return(nil, failure.Wrap(e)).
			Once()
		_, err := s.usecase.Find(ctx, model.UniqueDinosaurID(errUniqueID))
		s.True(errors.Is(err, e))
	}
}

func (s *UniqueDinosaurTestSuite) TestList() {
	uniques := model.UniqueDinosaurs{
		s.unique,
	}

	{
		s.mockDB.On(
			listUnique,
			ctx,
		).
			Return(uniques, nil).
			Once()
		r, err := s.usecase.List(ctx)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(uniques, r)
	}
	{
		s.mockDB.On(
			listUnique,
			ctx,
		).
			Return(nil, service.IntervalServerError).
			Once()
		_, err := s.usecase.List(ctx)
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDB.On(
			listUnique,
			ctx,
		).
			Return(nil, failure.Wrap(e)).
			Once()
		_, err := s.usecase.List(ctx)
		s.True(errors.Is(err, e))
	}
}

func (s *UniqueDinosaurTestSuite) TestInsert() {
	{
		s.mockDB.On(
			insertUnique,
			ctx,
			s.create,
		).
			Return(&s.unique, nil).
			Once()
		r, err := s.usecase.Create(ctx, s.create)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(&s.unique, r)
	}
	{
		s.mockDB.On(
			insertUnique,
			ctx,
			s.create,
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Create(ctx, s.create)
		s.True(errors.Is(err, e))
	}
}

func (s *UniqueDinosaurTestSuite) TestUpdate() {
	id := s.update.ID()
	{
		s.mockDB.On(findUnique, ctx, id).Return(&s.unique, nil).Once()
		s.mockDB.On(
			updateUnique,
			ctx,
			s.update,
		).
			Return(&s.unique, nil).
			Once()
		r, err := s.usecase.Update(ctx, s.update)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(&s.unique, r)
	}
	{
		s.mockDB.On(findUnique, ctx, id).Return(&s.unique, nil).Once()
		s.mockDB.On(
			updateUnique,
			ctx,
			s.update,
		).
			Return(nil, service.IntervalServerError).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDB.On(findUnique, ctx, id).Return(&s.unique, nil).Once()
		s.mockDB.On(
			updateUnique,
			ctx,
			s.update,
		).
			Return(nil, failure.Wrap(e)).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(errors.Is(err, e))
	}
}

func (s *UniqueDinosaurTestSuite) TestDelete() {
	{
		id := model.UniqueDinosaurID(uniqueID)
		s.mockDB.On(findUnique, ctx, id).Return(&s.unique, nil).Once()
		s.mockDB.On("Delete", ctx, id).Return(nil).Once()
		s.Nil(s.usecase.Delete(ctx, id))
	}
}

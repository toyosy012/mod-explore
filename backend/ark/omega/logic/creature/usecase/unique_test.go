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

	mockDinoCommand     *mockDinoCommandRepo
	mockUniqueQuery     *mockUniqueQueryRepo
	mockUniqueCommand   *mockUniqueCommandRepo
	mockVariantsCommand *mockVariantsCommandRepo
	usecase             UniqueUsecase

	create service.CreateCreature
	update service.UpdateCreature

	dinoResponse     service.ResponseDinosaur
	uniqueResponse   service.ResponseUnique
	variantsResponse service.ResponseVariants

	unique model.UniqueDinosaur
}

func newTestUniqueDinosaurSuite() *UniqueDinosaurTestSuite { return &UniqueDinosaurTestSuite{} }

func TestUniqueDinosaurSuite(t *testing.T) {
	suite.Run(t, newTestUniqueDinosaurSuite())
}

const (
	find   = "Select"
	list   = "List"
	insert = "Insert"
	update = "Update"
)

func (s *UniqueDinosaurTestSuite) SetupSuite() {
	{
		injector := do.New()
		mockDinoCommand := newMockDinoCommandRepo()
		do.ProvideValue[service.DinosaurCommandRepository](injector, mockDinoCommand)
		s.mockDinoCommand = mockDinoCommand
		mockUniqueQuery := newMockUniqueQuery()
		do.ProvideValue[service.UniqueQueryRepository](injector, mockUniqueQuery)
		s.mockUniqueQuery = mockUniqueQuery
		mockUniqueCommand := newMockUniqueCommand()
		do.ProvideValue[service.UniqueCommandRepository](injector, mockUniqueCommand)
		s.mockUniqueCommand = mockUniqueCommand
		mockVariantsCommand := newMockVariantsCommand()
		do.ProvideValue[service.VariantsCommandRepository](injector, mockVariantsCommand)
		s.mockVariantsCommand = mockVariantsCommand

		usecase, err := NewUnique(injector)
		if err != nil {
			return
		}
		s.usecase = usecase
	}
	{
		h, err := model.NewHealth(health)
		if err != nil {
			s.T().Error(err)
			return
		}
		m := model.NewMelee(melee)
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
				model.NewDinosaur(creatureID, creatureName, h, m),
				uniqueID, uniqueName, variants, *healthMultiplier, *meleeMultiplier,
			)
		}
		{
			s.create = service.NewCreateCreature(
				service.NewCreateDinosaur(creatureName, h, m),
				service.NewCreateUniqueDinosaur(uniqueName, *healthMultiplier, *meleeMultiplier),
				service.NewCreateVariants(uniqueID, variants),
			)
		}
		{
			s.update = service.NewUpdateCreature(
				service.NewUpdateDinosaur(creatureID, creatureName, health, melee),
				service.NewUpdateUniqueDinosaur(uniqueID, uniqueName, *healthMultiplier, *meleeMultiplier),
				service.NewUpdateVariants(variantsID, uniqueID, variants),
			)
		}
		{
			s.dinoResponse = service.NewResponseDinosaur(creatureID, creatureName, h, m)
			s.uniqueResponse = service.NewResponseUnique(uniqueID, uniqueName, *healthMultiplier, *meleeMultiplier)
			s.variantsResponse = service.NewResponseVariants(variantsID, uniqueID, variants)
		}
	}
}

const (
	multiplierHealth = 36.0
	health           = 1
	multiplierMelee  = 36.0
	melee            = 0
	cosmicID         = iota
	natureID
	creatureID   = 0
	creatureName = "dodo"
	uniqueID     = 0
	uniqueName   = "Kenny"
	cosmic       = "cosmic"
	singularity  = "singularity"
	nature       = "nature"
	thunderstorm = "thunderstorm"
	variantsID   = 0
)

func (s *UniqueDinosaurTestSuite) TestFind() {
	{
		s.mockUniqueQuery.On(
			find,
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
		s.mockUniqueQuery.On(
			find,
			ctx,
			model.UniqueDinosaurID(notExistUniqueID),
		).
			Return(nil, service.NotFound).
			Once()
		_, err := s.usecase.Find(ctx, model.UniqueDinosaurID(notExistUniqueID))
		s.True(failure.Is(err, logic.NotFound))
	}
	{
		s.mockUniqueQuery.On(
			find,
			ctx,
			model.UniqueDinosaurID(internalServerErrUniqueID),
		).
			Return(nil, service.IntervalServerError).
			Once()
		_, err := s.usecase.Find(ctx, model.UniqueDinosaurID(internalServerErrUniqueID))
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockUniqueQuery.On(
			find,
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
		s.mockUniqueQuery.On(
			list,
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
		s.mockUniqueQuery.On(
			list,
			ctx,
		).
			Return(nil, service.IntervalServerError).
			Once()
		_, err := s.usecase.List(ctx)
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockUniqueQuery.On(
			list,
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
		s.mockDinoCommand.On(
			insert,
			ctx,
			s.create.CreateDinosaur,
		).
			Return(&s.dinoResponse, nil).
			Once()
		s.mockUniqueCommand.On(
			insert,
			ctx,
			s.create.CreateUniqueDinosaur,
		).
			Return(&s.uniqueResponse, nil).
			Once()
		s.mockVariantsCommand.On(
			insert,
			ctx,
			s.create.CreateVariants,
		).
			Return(&s.variantsResponse, nil).
			Once()

		r, err := s.usecase.Create(ctx, s.create)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(&s.unique, r)
	}
	{
		s.mockDinoCommand.On(
			insert,
			ctx,
			s.create.CreateDinosaur,
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Create(ctx, s.create)
		s.True(errors.Is(err, e))
	}
	{
		s.mockDinoCommand.On(
			insert,
			ctx,
			s.create.CreateDinosaur,
		).
			Return(&s.dinoResponse, nil).
			Once()
		s.mockUniqueCommand.On(
			insert,
			ctx,
			s.create.CreateUniqueDinosaur,
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Create(ctx, s.create)
		s.True(errors.Is(err, e))
	}
	{
		s.mockDinoCommand.On(
			insert,
			ctx,
			s.create.CreateDinosaur,
		).
			Return(&s.dinoResponse, nil).
			Once()
		s.mockUniqueCommand.On(
			insert,
			ctx,
			s.create.CreateUniqueDinosaur,
		).
			Return(&s.uniqueResponse, nil).
			Once()
		s.mockVariantsCommand.On(
			insert,
			ctx,
			s.create.CreateVariants,
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Create(ctx, s.create)
		s.True(errors.Is(err, e))
	}
}

func (s *UniqueDinosaurTestSuite) TestUpdate() {
	id := s.update.ID()
	s.mockUniqueQuery.On(find, ctx, id).Return(&s.unique, nil).Times(7)
	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.UpdateDinosaur,
		).
			Return(&s.dinoResponse, nil).
			Once()
		s.mockUniqueCommand.On(
			update,
			ctx,
			s.update.UpdateUniqueDinosaur,
		).
			Return(&s.uniqueResponse, nil).
			Once()
		s.mockVariantsCommand.On(
			update,
			ctx,
			s.update.UpdateVariants,
		).
			Return(&s.variantsResponse, nil).
			Once()
		r, err := s.usecase.Update(ctx, s.update)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(&s.unique, r)
	}
	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.UpdateDinosaur,
		).
			Return(nil, service.IntervalServerError).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.UpdateDinosaur,
		).
			Return(&s.dinoResponse, nil).
			Once()
		s.mockUniqueCommand.On(
			update,
			ctx,
			s.update.UpdateUniqueDinosaur,
		).
			Return(nil, service.IntervalServerError).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.UpdateDinosaur,
		).
			Return(&s.dinoResponse, nil).
			Once()
		s.mockUniqueCommand.On(
			update,
			ctx,
			s.update.UpdateUniqueDinosaur,
		).
			Return(&s.uniqueResponse, nil).
			Once()
		s.mockVariantsCommand.On(
			update,
			ctx,
			s.update.UpdateVariants,
		).
			Return(nil, service.IntervalServerError).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(failure.Is(err, logic.IntervalServerError))
	}

	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.UpdateDinosaur,
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(errors.Is(err, e))
	}
	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.UpdateDinosaur,
		).
			Return(&s.dinoResponse, nil).
			Once()
		s.mockUniqueCommand.On(
			update,
			ctx,
			s.update.UpdateUniqueDinosaur,
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(errors.Is(err, e))
	}
	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.UpdateDinosaur,
		).
			Return(&s.dinoResponse, nil).
			Once()
		s.mockUniqueCommand.On(
			update,
			ctx,
			s.update.UpdateUniqueDinosaur,
		).
			Return(&s.uniqueResponse, nil).
			Once()
		s.mockVariantsCommand.On(
			update,
			ctx,
			s.update.UpdateVariants,
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(errors.Is(err, e))
	}
}

func (s *UniqueDinosaurTestSuite) TestDelete() {
	id := model.UniqueDinosaurID(uniqueID)
	s.mockUniqueQuery.On(find, ctx, id).Return(&s.unique, nil).Times(3)
	{
		s.mockUniqueCommand.On("Delete", ctx, id).Return(nil).Once()
		s.Nil(s.usecase.Delete(ctx, id))
	}
	{
		s.mockUniqueCommand.On("Delete", ctx, id).Return(service.IntervalServerError).Once()
		s.True(failure.Is(s.usecase.Delete(ctx, id), logic.IntervalServerError))
	}
	{
		s.mockUniqueCommand.On("Delete", ctx, id).Return(e).Once()
		s.True(errors.Is(s.usecase.Delete(ctx, id), e))
	}
}

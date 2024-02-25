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

	response service.ResponseCreature
	unique   model.UniqueDinosaur
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
		do.ProvideValue[service.UniqueVariantsCommand](injector, mockVariantsCommand)
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
			s.create = service.NewCreateCreature(
				creatureName, h, m, uniqueName, *healthMultiplier, *meleeMultiplier, [2]variantModel.VariantID{cosmicID, natureID},
			)
		}
		{
			s.update = service.NewUpdateCreature(
				creatureID, creatureName, health, melee,
				uniqueID, uniqueName, *healthMultiplier, *meleeMultiplier,
				variantsID, [2]variantModel.VariantID{variantsID},
			)
		}
		{
			s.dinoResponse = service.NewResponseDinosaur(creatureID, creatureName, h, m)
			s.uniqueResponse = service.NewResponseUnique(uniqueID, uniqueName, *healthMultiplier, *meleeMultiplier)
			s.variantsResponse = service.NewResponseVariants(variantsID, variants)

			s.response = service.ResponseCreature{
				ResponseDinosaur: s.dinoResponse,
				ResponseVariants: s.variantsResponse,
				ResponseUnique:   s.uniqueResponse,
			}
			s.unique = model.NewUniqueDinosaur(
				model.NewDinosaur(creatureID, creatureName, h, m),
				uniqueID, uniqueName,
				*healthMultiplier, *meleeMultiplier,
				variants,
			)
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
			Return(&s.response, nil).
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
	{
		s.mockUniqueQuery.On(
			list,
			ctx,
		).
			Return(service.ResponseCreatures{s.response}, nil).
			Once()
		r, err := s.usecase.List(ctx)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(model.UniqueDinosaurs{s.unique}, r)
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
			s.create.Dino(),
		).
			Return(model.DinosaurID(creatureID), nil).
			Once()
		s.mockVariantsCommand.On(
			insert,
			ctx,
			s.create.UniqueVariants(),
		).
			Return(model.UniqueVariantID(variantsID), nil).
			Once()
		s.mockUniqueCommand.On(
			insert,
			ctx,
			s.create.UniqueDinosaur(creatureID),
		).
			Return(model.UniqueDinosaurID(uniqueID), nil).
			Once()

		s.mockUniqueQuery.On(
			find,
			ctx,
			model.UniqueDinosaurID(uniqueID),
		).
			Return(&s.response, nil).
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
			s.create.Dino(),
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
			s.create.Dino(),
		).
			Return(model.DinosaurID(creatureID), nil).
			Once()
		s.mockUniqueCommand.On(
			insert,
			ctx,
			s.create.UniqueDinosaur(creatureID),
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
			s.create.Dino(),
		).
			Return(model.DinosaurID(creatureID), nil).
			Once()
		s.mockUniqueCommand.On(
			insert,
			ctx,
			s.create.UniqueDinosaur(creatureID),
		).
			Return(model.UniqueDinosaurID(uniqueID), nil).
			Once()
		s.mockVariantsCommand.On(
			insert,
			ctx,
			s.create.UniqueVariants(),
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Create(ctx, s.create)
		s.True(errors.Is(err, e))
	}
}

func (s *UniqueDinosaurTestSuite) TestUpdate() {
	id := s.update.Unique().ID()
	s.mockUniqueQuery.On(find, ctx, id).Return(&s.response, nil).Times(7)
	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.Dino(),
		).
			Return(nil).
			Once()
		s.mockVariantsCommand.On(
			update,
			ctx,
			s.update.Variants(),
		).
			Return(nil).
			Once()
		s.mockUniqueCommand.On(
			update,
			ctx,
			s.update.Unique(),
		).
			Return(nil).
			Once()
		s.mockUniqueQuery.On(find, ctx, id).Return(&s.response, nil).Once()
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
			s.update.Dino(),
		).
			Return(service.IntervalServerError).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.Dino(),
		).
			Return(nil).
			Once()
		s.mockUniqueCommand.On(
			update,
			ctx,
			s.update.Unique(),
		).
			Return(service.IntervalServerError).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.Dino(),
		).
			Return(nil).
			Once()
		s.mockUniqueCommand.On(
			update,
			ctx,
			s.update.Unique(),
		).
			Return(nil).
			Once()
		s.mockVariantsCommand.On(
			update,
			ctx,
			s.update.Variants(),
		).
			Return(service.IntervalServerError).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(failure.Is(err, logic.IntervalServerError))
	}

	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.Dino(),
		).
			Return(e).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(errors.Is(err, e))
	}
	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.Dino(),
		).
			Return(nil).
			Once()
		s.mockUniqueCommand.On(
			update,
			ctx,
			s.update.Unique(),
		).
			Return(e).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(errors.Is(err, e))
	}
	{
		s.mockDinoCommand.On(
			update,
			ctx,
			s.update.Dino(),
		).
			Return(nil).
			Once()
		s.mockUniqueCommand.On(
			update,
			ctx,
			s.update.Unique(),
		).
			Return(nil).
			Once()
		s.mockVariantsCommand.On(
			update,
			ctx,
			s.update.Variants(),
		).
			Return(e).
			Once()
		_, err := s.usecase.Update(ctx, s.update)
		s.True(errors.Is(err, e))
	}
}

func (s *UniqueDinosaurTestSuite) TestDelete() {
	id := model.UniqueDinosaurID(uniqueID)
	s.mockUniqueQuery.On(find, ctx, id).Return(&s.response, nil).Times(3)
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

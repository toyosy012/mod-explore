package usecase

import (
	"errors"
	"testing"

	"github.com/morikuni/failure"
	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"mods-explore/ark/omega/logic"
	"mods-explore/ark/omega/logic/variant/domain/model"
	"mods-explore/ark/omega/logic/variant/domain/service"
)

type VariantTestSuite struct {
	suite.Suite

	mockDB  *mockDBClient
	usecase VariantUsecase
}

func newTestVariantSuite() *VariantTestSuite { return &VariantTestSuite{} }

func TestVariantSuite(t *testing.T) {
	suite.Run(t, newTestVariantSuite())
}

const (
	findVariant   = "FindVariant"
	listVariant   = "ListVariants"
	createVariant = "CreateVariant"
	updateVariant = "UpdateVariant"
	deleteVariant = "DeleteVariant"
)

func (s *VariantTestSuite) SetupSuite() {
	injector := do.New()

	mockDB := newMockDBClient()
	do.ProvideValue[service.VariantRepository](injector, mockDB)
	s.mockDB = mockDB
	usecase, err := NewVariant(injector)
	if err != nil {
		return
	}

	s.usecase = usecase
}

func (s *VariantTestSuite) TestFind() {
	variant := model.NewVariant(model.VariantID(id), "cosmic", "meteor")
	{
		s.mockDB.On(
			findVariant,
			ctx,
			model.VariantID(id),
		).
			Return(&variant, nil)
		r, err := s.usecase.Find(ctx, model.VariantID(id))
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(&variant, r)
	}
	{
		s.mockDB.On(
			findVariant,
			ctx,
			model.VariantID(notExistID),
		).
			Return(nil, service.NotFound)
		_, err := s.usecase.Find(ctx, model.VariantID(notExistID))
		s.True(failure.Is(err, logic.NotFound))
	}
	{
		s.mockDB.On(
			findVariant,
			ctx,
			model.VariantID(intervalServerErrID),
		).
			Return(nil, service.IntervalServerError)
		_, err := s.usecase.Find(ctx, model.VariantID(intervalServerErrID))
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDB.On(
			findVariant,
			ctx,
			model.VariantID(errID),
		).
			Return(nil, failure.Wrap(e))
		_, err := s.usecase.Find(ctx, model.VariantID(errID))
		s.True(errors.Is(err, e))
	}
}

func (s *VariantTestSuite) TestList() {
	variants := model.Variants{
		model.NewVariant(model.VariantID(id), "cosmic", "meteor"),
	}

	{
		s.mockDB.On(
			listVariant,
			ctx,
		).
			Return(variants, nil).
			Once()
		r, err := s.usecase.List(ctx)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(variants, r)
	}
	{
		s.mockDB.On(
			listVariant,
			ctx,
		).
			Return(nil, service.IntervalServerError).
			Once()
		_, err := s.usecase.List(ctx)
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDB.On(
			listVariant,
			ctx,
		).
			Return(nil, failure.Wrap(e)).
			Once()
		_, err := s.usecase.List(ctx)
		s.True(errors.Is(err, e))
	}
}

func (s *VariantTestSuite) TestCreate() {
	item := service.NewCreateVariant(groupID, "meteor")
	variant := model.NewVariant(id, "cosmic", "meteor")
	{
		s.mockDB.On(
			createVariant,
			ctx,
			item,
		).
			Return(&variant, nil).
			Once()
		r, err := s.usecase.Create(ctx, item)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(&variant, r)
	}
	{
		s.mockDB.On(
			createVariant,
			ctx,
			item,
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Create(ctx, item)
		s.True(errors.Is(err, e))
	}
}

// TestUpdate レコードの存在確認はFindと同じなのチェックしない
func (s *VariantTestSuite) TestUpdate() {
	{
		item := service.NewUpdateVariant(id, groupID, "meteor")
		variant := model.NewVariant(id, "cosmic", "meteor")
		s.mockDB.On(findVariant, model.VariantID(id)).Return(&variant, nil).Once()
		s.mockDB.On(
			updateVariant,
			ctx,
			item,
		).
			Return(&variant, nil).
			Once()
		r, err := s.usecase.Update(ctx, item)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(&variant, r)
	}

	{ // updateVariant error case
		item := service.NewUpdateVariant(id, groupID, "meteor")
		variant := model.NewVariant(id, "cosmic", "meteor")
		s.mockDB.On(findVariant, model.VariantID(id)).Return(&variant, nil).Once()
		s.mockDB.On(
			updateVariant,
			ctx,
			item,
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Update(ctx, item)
		s.True(errors.Is(err, e))
	}
}

func (s *VariantTestSuite) TestDelete() {
	{
		variant := model.NewVariant(id, "cosmic", "meteor")
		s.mockDB.On(findVariant, ctx, model.VariantID(id)).Return(&variant, nil).Once()
		s.mockDB.On(
			deleteVariant,
			ctx,
			model.VariantID(id),
		).
			Return(nil).
			Once()
		s.Nil(s.usecase.Delete(ctx, model.VariantID(id)))
	}

	{ // updateVariant error case
		variant := model.NewVariant(id, "cosmic", "meteor")
		s.mockDB.On(findVariant, ctx, model.VariantID(id)).Return(&variant, nil).Once()
		s.mockDB.On(
			deleteVariant,
			ctx,
			model.VariantID(id),
		).
			Return(e).
			Once()
		err := s.usecase.Delete(ctx, model.VariantID(id))
		s.True(errors.Is(err, e))
	}
}

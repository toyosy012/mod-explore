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

type VariantGroupTestSuite struct {
	suite.Suite

	mockDB  *mockVariantGroup
	usecase VariantGroupUsecase
}

func newTestVariantGroupSuite() *VariantGroupTestSuite { return &VariantGroupTestSuite{} }

func TestVariantGroupSuite(t *testing.T) {
	suite.Run(t, newTestVariantGroupSuite())
}

const (
	findVariantGroup   = "Select"
	listVariantGroup   = "List"
	createVariantGroup = "Insert"
	updateVariantGroup = "Update"
	deleteVariantGroup = "Delete"
)

func (s *VariantGroupTestSuite) SetupSuite() {
	injector := do.New()

	mockDB := newMockVariantGroup()
	do.ProvideValue[service.VariantGroupRepository](injector, mockDB)
	s.mockDB = mockDB
	usecase, err := NewVariantGroup(injector)
	if err != nil {
		return
	}

	s.usecase = usecase
}

func (s *VariantGroupTestSuite) TestFind() {
	variantGroup := model.NewVariantGroup(model.VariantGroupID(groupID), "cosmic")
	{
		s.mockDB.On(
			findVariantGroup,
			ctx,
			model.VariantGroupID(groupID),
		).
			Return(&variantGroup, nil)
		r, err := s.usecase.Find(ctx, model.VariantGroupID(groupID))
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(&variantGroup, r)
	}
	{
		s.mockDB.On(
			findVariantGroup,
			ctx,
			model.VariantGroupID(notExistGroupID),
		).
			Return(nil, service.NotFound)
		_, err := s.usecase.Find(ctx, model.VariantGroupID(notExistGroupID))
		s.True(failure.Is(err, logic.NotFound))
	}
	{
		s.mockDB.On(
			findVariantGroup,
			ctx,
			model.VariantGroupID(internalServerErrGroupID),
		).
			Return(nil, service.IntervalServerError)
		_, err := s.usecase.Find(ctx, model.VariantGroupID(internalServerErrGroupID))
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDB.On(
			findVariantGroup,
			ctx,
			model.VariantGroupID(errGroupID),
		).
			Return(nil, failure.Wrap(e))
		_, err := s.usecase.Find(ctx, model.VariantGroupID(errGroupID))
		s.True(errors.Is(err, e))
	}
}

func (s *VariantGroupTestSuite) TestList() {
	variantGroups := model.VariantGroups{
		model.NewVariantGroup(model.VariantGroupID(groupID), "cosmic"),
	}

	{
		s.mockDB.On(
			listVariantGroup,
			ctx,
		).
			Return(variantGroups, nil).
			Once()
		r, err := s.usecase.List(ctx)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(variantGroups, r)
	}
	{
		s.mockDB.On(
			listVariantGroup,
			ctx,
		).
			Return(nil, service.IntervalServerError).
			Once()
		_, err := s.usecase.List(ctx)
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDB.On(
			listVariantGroup,
			ctx,
		).
			Return(nil, failure.Wrap(e)).
			Once()
		_, err := s.usecase.List(ctx)
		s.True(errors.Is(err, e))
	}
}

func (s *VariantGroupTestSuite) TestCreate() {
	item := service.NewCreateVariantGroup("cosmic")
	variantGroup := model.NewVariantGroup(groupID, "cosmic")
	{
		s.mockDB.On(
			createVariantGroup,
			ctx,
			item,
		).
			Return(&variantGroup, nil).
			Once()
		r, err := s.usecase.Create(ctx, item)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(&variantGroup, r)
	}
	{
		s.mockDB.On(
			createVariantGroup,
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
func (s *VariantGroupTestSuite) TestUpdate() {
	{
		item := service.NewUpdateVariantGroup(groupID, "cosmic")
		variantGroup := model.NewVariantGroup(groupID, "cosmic")
		s.mockDB.On(findVariantGroup, model.VariantGroupID(groupID)).Return(&variantGroup, nil).Once()
		s.mockDB.On(
			updateVariantGroup,
			ctx,
			item,
		).
			Return(&variantGroup, nil).
			Once()
		r, err := s.usecase.Update(ctx, item)
		if err != nil {
			s.T().Error(err)
			return
		}

		s.Equal(&variantGroup, r)
	}

	{ // updateVariant error case
		item := service.NewUpdateVariantGroup(groupID, "cosmic")
		variantGroup := model.NewVariantGroup(groupID, "cosmic")
		s.mockDB.On(findVariantGroup, model.VariantGroupID(groupID)).Return(&variantGroup, nil).Once()
		s.mockDB.On(
			updateVariantGroup,
			ctx,
			item,
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Update(ctx, item)
		s.True(errors.Is(err, e))
	}
}

func (s *VariantGroupTestSuite) TestDelete() {
	{
		variant := model.NewVariantGroup(groupID, "cosmic")
		s.mockDB.On(findVariantGroup, ctx, model.VariantGroupID(groupID)).Return(&variant, nil).Once()
		s.mockDB.On(
			deleteVariantGroup,
			ctx,
			model.VariantGroupID(groupID),
		).
			Return(nil).
			Once()
		s.Nil(s.usecase.Delete(ctx, model.VariantGroupID(groupID)))
	}

	{ // updateVariant error case
		variant := model.NewVariantGroup(groupID, "cosmic")
		s.mockDB.On(findVariantGroup, ctx, model.VariantGroupID(groupID)).Return(&variant, nil).Once()
		s.mockDB.On(
			deleteVariantGroup,
			ctx,
			model.VariantGroupID(groupID),
		).
			Return(e).
			Once()
		err := s.usecase.Delete(ctx, model.VariantGroupID(groupID))
		s.True(errors.Is(err, e))
	}
}

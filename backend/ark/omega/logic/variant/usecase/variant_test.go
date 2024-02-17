package usecase

import (
	"context"
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

var (
	ctx    = context.Background()
	find   = "FindVariant"
	create = "CreateVariant"
	e      = errors.New("test")
)

const (
	id = iota
	notExistID
	intervalServerErrID
	errID
)

const (
	groupID = iota
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
			find,
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
			find,
			ctx,
			model.VariantID(notExistID),
		).
			Return(nil, service.NotFound)
		_, err := s.usecase.Find(ctx, model.VariantID(notExistID))
		s.True(failure.Is(err, logic.NotFound))
	}
	{
		s.mockDB.On(
			find,
			ctx,
			model.VariantID(intervalServerErrID),
		).
			Return(nil, service.IntervalServerError)
		_, err := s.usecase.Find(ctx, model.VariantID(intervalServerErrID))
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDB.On(
			find,
			ctx,
			model.VariantID(errID),
		).
			Return(nil, failure.Wrap(e))
		_, err := s.usecase.Find(ctx, model.VariantID(errID))
		s.True(errors.Is(err, e))
	}
}

func (s *VariantTestSuite) TestList() {
	method := "ListVariants"
	variants := model.Variants{
		model.NewVariant(model.VariantID(id), "cosmic", "meteor"),
	}

	{
		s.mockDB.On(
			method,
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
			method,
			ctx,
		).
			Return(nil, service.IntervalServerError).
			Once()
		_, err := s.usecase.List(ctx)
		s.True(failure.Is(err, logic.IntervalServerError))
	}
	{
		s.mockDB.On(
			method,
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
			create,
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
			create,
			ctx,
			item,
		).
			Return(nil, e).
			Once()
		_, err := s.usecase.Create(ctx, item)
		s.True(errors.Is(err, e))
	}
}

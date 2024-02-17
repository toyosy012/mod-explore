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
	ctx  = context.Background()
	find = "FindVariant"
	e    = errors.New("test")
)

const (
	id = iota
	notExistID
	intervalServerErrID
	errID
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

package usecase

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"mods-explore/ark/omega/logic"
	"mods-explore/ark/omega/logic/creature/domain/model"
	"mods-explore/ark/omega/logic/creature/domain/service"
)

var (
	ctx = context.Background()
	e   = errors.New("test")
)

var _ logic.Transactioner = (*mockUniqueQueryRepo)(nil)
var _ service.UniqueQueryRepository = (*mockUniqueQueryRepo)(nil)

type mockUniqueQueryRepo struct {
	mock.Mock
}

func newMockUniqueQuery() *mockUniqueQueryRepo { return &mockUniqueQueryRepo{} }

func (g *mockUniqueQueryRepo) WithTransaction(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
	return fn(ctx)
}

func (g *mockUniqueQueryRepo) Select(ctx context.Context, id model.UniqueDinosaurID) (*model.UniqueDinosaur, error) {
	args := g.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UniqueDinosaur), args.Error(1)
}

func (g *mockUniqueQueryRepo) List(ctx context.Context) (model.UniqueDinosaurs, error) {
	args := g.Called(ctx)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(model.UniqueDinosaurs), nil
}

var _ logic.Transactioner = (*mockUniqueCommandRepo)(nil)
var _ service.UniqueCommandRepository = (*mockUniqueCommandRepo)(nil)

type mockUniqueCommandRepo struct {
	mock.Mock
}

func newMockUniqueCommand() *mockUniqueCommandRepo { return &mockUniqueCommandRepo{} }

func (g *mockUniqueCommandRepo) WithTransaction(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
	return fn(ctx)
}

func (g *mockUniqueCommandRepo) Insert(
	ctx context.Context, create service.CreateUniqueDinosaur,
) (*model.UniqueDinosaur, error) {
	args := g.Called(ctx, create)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UniqueDinosaur), args.Error(1)
}

func (g *mockUniqueCommandRepo) Update(
	ctx context.Context, update service.UpdateUniqueDinosaur,
) (*model.UniqueDinosaur, error) {
	args := g.Called(ctx, update)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UniqueDinosaur), args.Error(1)
}

func (g *mockUniqueCommandRepo) Delete(ctx context.Context, id model.UniqueDinosaurID) error {
	args := g.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil
	}
	return args.Error(0)
}

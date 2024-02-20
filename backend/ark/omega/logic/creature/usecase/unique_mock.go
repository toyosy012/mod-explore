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

var _ logic.Transactioner = (*mockDinoCommandRepo)(nil)
var _ service.DinosaurCommandRepository = (*mockDinoCommandRepo)(nil)

type mockDinoCommandRepo struct {
	mock.Mock
}

func newMockDinoCommandRepo() *mockDinoCommandRepo { return &mockDinoCommandRepo{} }

func (g *mockDinoCommandRepo) WithTransaction(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
	return fn(ctx)
}

func (g *mockDinoCommandRepo) Insert(
	ctx context.Context, create service.CreateDinosaur,
) (model.DinosaurID, error) {
	args := g.Called(ctx, create)

	r := args.Get(0)
	if nil == r {
		return 0, args.Error(1)
	}
	return r.(model.DinosaurID), nil
}

func (g *mockDinoCommandRepo) Update(
	ctx context.Context, update service.UpdateDinosaur,
) error {
	args := g.Called(ctx, update)

	return args.Error(0)
}

func (g *mockDinoCommandRepo) Delete(ctx context.Context, id model.DinosaurID) error {
	args := g.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil
	}
	return args.Error(0)
}

var _ logic.Transactioner = (*mockUniqueQueryRepo)(nil)
var _ service.UniqueQueryRepository = (*mockUniqueQueryRepo)(nil)

type mockUniqueQueryRepo struct {
	mock.Mock
}

func newMockUniqueQuery() *mockUniqueQueryRepo { return &mockUniqueQueryRepo{} }

func (g *mockUniqueQueryRepo) WithTransaction(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
	return fn(ctx)
}

func (g *mockUniqueQueryRepo) Select(ctx context.Context, id model.UniqueDinosaurID) (*service.ResponseCreature, error) {
	args := g.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.ResponseCreature), args.Error(1)
}

func (g *mockUniqueQueryRepo) List(ctx context.Context) (service.ResponseCreatures, error) {
	args := g.Called(ctx)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(service.ResponseCreatures), nil
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
) (model.UniqueDinosaurID, error) {
	args := g.Called(ctx, create)

	r := args.Get(0)
	if r == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(model.UniqueDinosaurID), args.Error(1)
}

func (g *mockUniqueCommandRepo) Update(
	ctx context.Context, update service.UpdateUniqueDinosaur,
) error {
	args := g.Called(ctx, update)

	return args.Error(0)
}

func (g *mockUniqueCommandRepo) Delete(ctx context.Context, id model.UniqueDinosaurID) error {
	args := g.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil
	}
	return args.Error(0)
}

var _ logic.Transactioner = (*mockVariantsCommandRepo)(nil)
var _ service.VariantsCommandRepository = (*mockVariantsCommandRepo)(nil)

type mockVariantsCommandRepo struct {
	mock.Mock
}

func newMockVariantsCommand() *mockVariantsCommandRepo { return &mockVariantsCommandRepo{} }

func (e *mockVariantsCommandRepo) WithTransaction(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
	return fn(ctx)
}

func (e *mockVariantsCommandRepo) Insert(
	ctx context.Context, create service.CreateVariants,
) (model.UniqueVariantID, error) {
	args := e.Called(ctx, create)

	r := args.Get(0)
	if r == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(model.UniqueVariantID), args.Error(1)
}

func (e *mockVariantsCommandRepo) Update(
	ctx context.Context, update service.UpdateVariants,
) error {
	args := e.Called(ctx, update)

	return args.Error(0)
}

func (e *mockVariantsCommandRepo) Delete(ctx context.Context, id model.UniqueVariantID) error {
	args := e.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil
	}
	return args.Error(0)
}

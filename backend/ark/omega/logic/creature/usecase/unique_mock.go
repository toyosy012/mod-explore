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

var _ logic.Transactioner = (*mockUniqueDB)(nil)
var _ service.UniqueRepository = (*mockUniqueDB)(nil)

type mockUniqueDB struct {
	mock.Mock
}

func newMockUniqueDB() *mockUniqueDB { return &mockUniqueDB{} }

func (g *mockUniqueDB) WithTransaction(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
	return fn(ctx)
}

func (g *mockUniqueDB) Select(ctx context.Context, id model.UniqueDinosaurID) (*model.UniqueDinosaur, error) {
	args := g.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UniqueDinosaur), args.Error(1)
}

func (g *mockUniqueDB) List(ctx context.Context) (model.UniqueDinosaurs, error) {
	args := g.Called(ctx)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(model.UniqueDinosaurs), nil
}

func (g *mockUniqueDB) Insert(
	ctx context.Context, create service.CreateUniqueDinosaur,
) (*model.UniqueDinosaur, error) {
	args := g.Called(ctx, create)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UniqueDinosaur), args.Error(1)
}

func (g *mockUniqueDB) Update(
	ctx context.Context, update service.UpdateUniqueDinosaur,
) (*model.UniqueDinosaur, error) {
	args := g.Called(ctx, update)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UniqueDinosaur), args.Error(1)
}

func (g *mockUniqueDB) Delete(ctx context.Context, id model.UniqueDinosaurID) error {
	args := g.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil
	}
	return args.Error(0)
}

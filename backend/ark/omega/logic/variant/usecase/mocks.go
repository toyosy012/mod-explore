package usecase

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"mods-explore/ark/omega/logic"
	"mods-explore/ark/omega/logic/variant/domain/model"
	"mods-explore/ark/omega/logic/variant/domain/service"
)

var (
	ctx = context.Background()
	e   = errors.New("test")
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

// mockTransactionがインターフェースを満たしているか
var _ logic.Transactioner = (*mockDBClient)(nil)
var _ service.VariantRepository = (*mockDBClient)(nil)

type mockDBClient struct {
	mock.Mock
}

func newMockDBClient() *mockDBClient { return &mockDBClient{} }

func (c *mockDBClient) WithTransaction(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
	return fn(ctx)
}

func (c *mockDBClient) FindVariant(ctx context.Context, id model.VariantID) (*model.Variant, error) {
	args := c.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Variant), args.Error(1)
}

func (c *mockDBClient) ListVariants(ctx context.Context) (model.Variants, error) {
	args := c.Called(ctx)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(model.Variants), nil
}
func (c *mockDBClient) CreateVariant(ctx context.Context, create service.CreateVariant) (*model.Variant, error) {
	args := c.Called(ctx, create)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Variant), args.Error(1)
}
func (c *mockDBClient) UpdateVariant(ctx context.Context, update service.UpdateVariant) (*model.Variant, error) {
	args := c.Called(ctx, update)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Variant), args.Error(1)
}
func (c *mockDBClient) DeleteVariant(ctx context.Context, id model.VariantID) error {
	args := c.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil
	}
	return args.Error(0)
}

var _ logic.Transactioner = (*mockVariantGroup)(nil)
var _ service.VariantGroupRepository = (*mockVariantGroup)(nil)

type mockVariantGroup struct {
	mock.Mock
}

func newMockVariantGroup() *mockVariantGroup { return &mockVariantGroup{} }

func (g *mockVariantGroup) WithTransaction(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
	return fn(ctx)
}

func (g *mockVariantGroup) Select(ctx context.Context, id model.VariantGroupID) (*model.VariantGroup, error) {
	args := g.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.VariantGroup), args.Error(1)
}

func (g *mockVariantGroup) List(ctx context.Context) (model.VariantGroups, error) {
	args := g.Called(ctx)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(model.VariantGroups), nil
}

func (g *mockVariantGroup) Insert(ctx context.Context, create service.CreateVariantGroup) (*model.VariantGroup, error) {
	args := g.Called(ctx, create)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.VariantGroup), args.Error(1)
}

func (g *mockVariantGroup) Update(ctx context.Context, update service.UpdateVariantGroup) (*model.VariantGroup, error) {
	args := g.Called(ctx, update)

	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.VariantGroup), args.Error(1)
}

func (g *mockVariantGroup) Delete(ctx context.Context, id model.VariantGroupID) error {
	args := g.Called(ctx, id)

	r := args.Get(0)
	if r == nil {
		return nil
	}
	return args.Error(0)
}

package usecase

import (
	"context"
	"errors"

	"github.com/morikuni/failure"

	"mods-explore/ark/omega/logic"
	"mods-explore/ark/omega/logic/variant/domain/model"
	"mods-explore/ark/omega/logic/variant/domain/service"
)

type VariantUsecase interface {
	Find(context.Context, model.VariantID) (*model.Variant, error)
	List(context.Context) (model.Variants, error)
	Create(context.Context, service.CreateVariant) (*model.Variant, error)
	Update(context.Context, service.UpdateVariant) (*model.Variant, error)
	Delete(context.Context, model.VariantID) error
}

type Variant struct {
	repository service.VariantRepository
}

func NewVariant(repo service.VariantRepository) VariantUsecase {
	return Variant{
		repository: repo,
	}
}

func (v Variant) Find(ctx context.Context, id model.VariantID) (*model.Variant, error) {
	variant, err := v.repository.FindVariant(ctx, id)
	if err != nil {
		if errors.Is(err, service.NotFound) {
			return nil, failure.New(logic.NotFound)
		} else if errors.Is(err, service.IntervalServerError) {
			return nil, failure.New(logic.IntervalServerError)
		}
		return nil, failure.Wrap(err)
	}

	return variant, nil
}

func (v Variant) List(ctx context.Context) (model.Variants, error) {
	variants, err := v.repository.ListVariants(ctx)
	if err != nil {
		if errors.Is(err, service.IntervalServerError) {
			return nil, failure.New(logic.IntervalServerError)
		}
		return nil, failure.Wrap(err)
	}
	return variants, nil
}

func (v Variant) Create(ctx context.Context, item service.CreateVariant) (*model.Variant, error) {
	variant, err := v.repository.CreateVariant(ctx, item)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return variant, nil
}

func (v Variant) Update(ctx context.Context, item service.UpdateVariant) (*model.Variant, error) {
	variant, err := v.repository.UpdateVariant(ctx, item)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return variant, nil
}

func (v Variant) Delete(ctx context.Context, id model.VariantID) error {
	err := v.repository.DeleteVariant(ctx, id)
	if err != nil {
		if errors.Is(err, service.NotFound) {
			return failure.New(logic.NotFound)
		} else if errors.Is(err, service.IntervalServerError) {
			return failure.New(logic.IntervalServerError)
		}
		return failure.Wrap(err)
	}
	return nil
}

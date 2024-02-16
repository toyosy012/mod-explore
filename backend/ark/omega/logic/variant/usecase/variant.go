package usecase

import (
	"context"
	"errors"

	"github.com/morikuni/failure"
	"github.com/samber/do"

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

func NewVariant(injector *do.Injector) (VariantUsecase, error) {
	return Variant{
		repository: do.MustInvoke[service.VariantRepository](injector),
	}, nil
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
	return logic.UseTransactioner(ctx, func(ctx context.Context) (*model.Variant, error) {
		variant, err := v.repository.CreateVariant(ctx, item)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		return variant, nil
	})
}

func (v Variant) Update(ctx context.Context, item service.UpdateVariant) (*model.Variant, error) {
	return logic.UseTransactioner(ctx, func(ctx context.Context) (*model.Variant, error) {
		if _, err := v.repository.FindVariant(ctx, item.ID()); err != nil {
			if errors.Is(err, service.NotFound) {
				return nil, failure.New(logic.NotFound)
			}
			return nil, failure.Wrap(err)
		}

		variant, err := v.repository.UpdateVariant(ctx, item)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		return variant, nil
	})
}

func (v Variant) Delete(ctx context.Context, id model.VariantID) error {
	return logic.UseTransactioner0(ctx, func(ctx context.Context) error {
		_, err := v.repository.FindVariant(ctx, id)
		if errors.Is(err, service.NotFound) {
			return failure.New(logic.NotFound)
		}
		if err != nil {
			return err
		}

		err = v.repository.DeleteVariant(ctx, id)
		if err != nil {
			if errors.Is(err, service.NotFound) {
				return failure.New(logic.NotFound)
			} else if errors.Is(err, service.IntervalServerError) {
				return failure.New(logic.IntervalServerError)
			}
			return failure.Wrap(err)
		}
		return nil
	})
}

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

type VariantGroupUsecase interface {
	Find(context.Context, model.VariantGroupID) (*model.VariantGroup, error)
	List(context.Context) (model.VariantGroups, error)
	Create(context.Context, service.CreateVariantGroup) (*model.VariantGroup, error)
	Update(context.Context, service.UpdateVariantGroup) (*model.VariantGroup, error)
	Delete(context.Context, model.VariantGroupID) error
}

type VariantGroup struct {
	repository service.VariantGroupRepository
}

func NewVariantGroup(injector *do.Injector) (VariantGroupUsecase, error) {
	return &VariantGroup{
		repository: do.MustInvoke[service.VariantGroupRepository](injector),
	}, nil
}

func (v VariantGroup) Find(ctx context.Context, id model.VariantGroupID) (*model.VariantGroup, error) {
	variant, err := v.repository.Select(ctx, id)
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

func (v VariantGroup) List(ctx context.Context) (model.VariantGroups, error) {
	variants, err := v.repository.List(ctx)
	if err != nil {
		if errors.Is(err, service.IntervalServerError) {
			return nil, failure.New(logic.IntervalServerError)
		}
		return nil, failure.Wrap(err)
	}
	return variants, nil
}

func (v VariantGroup) Create(ctx context.Context, item service.CreateVariantGroup) (*model.VariantGroup, error) {
	return logic.UseTransactioner(ctx, func(ctx context.Context) (*model.VariantGroup, error) {
		variant, err := v.repository.Insert(ctx, item)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		return variant, nil
	})
}

func (v VariantGroup) Update(ctx context.Context, item service.UpdateVariantGroup) (*model.VariantGroup, error) {
	return logic.UseTransactioner(ctx, func(ctx context.Context) (*model.VariantGroup, error) {
		if _, err := v.repository.Select(ctx, item.ID()); err != nil {
			if errors.Is(err, service.NotFound) {
				return nil, failure.New(logic.NotFound)
			}
			return nil, failure.Wrap(err)
		}

		variant, err := v.repository.Update(ctx, item)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		return variant, nil
	})
}

func (v VariantGroup) Delete(ctx context.Context, id model.VariantGroupID) error {
	return logic.UseTransactioner0(ctx, func(ctx context.Context) error {
		_, err := v.repository.Select(ctx, id)
		if errors.Is(err, service.NotFound) {
			return failure.New(logic.NotFound)
		}
		if err != nil {
			return err
		}

		err = v.repository.Delete(ctx, id)
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

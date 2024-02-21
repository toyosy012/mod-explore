package usecase

import (
	"context"
	"errors"

	"github.com/morikuni/failure"
	"github.com/samber/do"
	"github.com/samber/lo"

	"mods-explore/ark/omega/logic"
	"mods-explore/ark/omega/logic/creature/domain/model"
	"mods-explore/ark/omega/logic/creature/domain/service"
)

type UniqueUsecase interface {
	Find(context.Context, model.UniqueDinosaurID) (*model.UniqueDinosaur, error)
	List(context.Context) (model.UniqueDinosaurs, error)
	Create(context.Context, service.CreateCreature) (*model.UniqueDinosaur, error)
	Update(context.Context, service.UpdateCreature) (*model.UniqueDinosaur, error)
	Delete(context.Context, model.UniqueDinosaurID) error
}

type Unique struct {
	dinoCommand    service.DinosaurCommandRepository
	uniqueQuery    service.UniqueQueryRepository
	uniqueCommand  service.UniqueCommandRepository
	variantCommand service.UniqueVariantsCommand
}

func NewUnique(injector *do.Injector) (*Unique, error) {
	return &Unique{
		dinoCommand:    do.MustInvoke[service.DinosaurCommandRepository](injector),
		uniqueQuery:    do.MustInvoke[service.UniqueQueryRepository](injector),
		uniqueCommand:  do.MustInvoke[service.UniqueCommandRepository](injector),
		variantCommand: do.MustInvoke[service.UniqueVariantsCommand](injector),
	}, nil
}

func (u Unique) Find(ctx context.Context, id model.UniqueDinosaurID) (*model.UniqueDinosaur, error) {
	resp, err := u.uniqueQuery.Select(ctx, id)
	if err != nil {
		if errors.Is(err, service.NotFound) {
			return nil, failure.New(logic.NotFound)
		} else if errors.Is(err, service.IntervalServerError) {
			return nil, failure.New(logic.IntervalServerError)
		}
		return nil, failure.Wrap(err)
	}

	unique := resp.ToUniqueDinosaur()
	return &unique, nil
}

func (u Unique) List(ctx context.Context) (model.UniqueDinosaurs, error) {
	resp, err := u.uniqueQuery.List(ctx)
	if err != nil {
		if errors.Is(err, service.IntervalServerError) {
			return nil, failure.New(logic.IntervalServerError)
		}
		return nil, failure.Wrap(err)
	}

	uniques := lo.Map(resp, func(r service.ResponseCreature, _ int) model.UniqueDinosaur {
		return r.ToUniqueDinosaur()
	})
	return uniques, nil
}

func (u Unique) Create(ctx context.Context, create service.CreateCreature) (_ *model.UniqueDinosaur, err error) {
	return logic.UseTransactioner(ctx, func(ctx context.Context) (*model.UniqueDinosaur, error) {
		var dinoID model.DinosaurID
		if dinoID, err = u.dinoCommand.Insert(
			ctx,
			create.Dino(),
		); err != nil {
			return nil, failure.Wrap(err)
		}
		var uniqueVariantID model.UniqueVariantID
		if uniqueVariantID, err = u.variantCommand.Insert(
			ctx,
			create.UniqueVariants(),
		); err != nil {
			return nil, failure.Wrap(err)
		}
		var uniqueID model.UniqueDinosaurID
		if uniqueID, err = u.uniqueCommand.Insert(
			ctx,
			service.NewCreateUniqueDinosaur(
				create.UniqueName, create.HealthMultiplier, create.DamageMultiplier, dinoID, uniqueVariantID,
			),
		); err != nil {
			return nil, failure.Wrap(err)
		}

		resp, err := u.uniqueQuery.Select(ctx, uniqueID)
		if err != nil {
			if errors.Is(err, service.NotFound) {
				return nil, failure.New(logic.NotFound)
			}
			if errors.Is(err, service.IntervalServerError) {
				return nil, failure.New(logic.IntervalServerError)
			}
			return nil, failure.Wrap(err)
		}

		unique := resp.ToUniqueDinosaur()
		return &unique, nil
	})
}

func (u Unique) Update(ctx context.Context, update service.UpdateCreature) (_ *model.UniqueDinosaur, err error) {
	return logic.UseTransactioner(ctx, func(ctx context.Context) (*model.UniqueDinosaur, error) {
		if _, err = u.uniqueQuery.Select(ctx, update.Unique().ID()); err != nil {
			if errors.Is(err, service.NotFound) {
				return nil, failure.New(logic.NotFound)
			}
			return nil, failure.Wrap(err)
		}

		if err = u.dinoCommand.Update(ctx, update.Dino()); err != nil {
			if errors.Is(err, service.IntervalServerError) {
				return nil, failure.New(logic.IntervalServerError)
			}
			return nil, failure.Wrap(err)
		}

		if err = u.uniqueCommand.Update(ctx, update.Unique()); err != nil {
			if errors.Is(err, service.IntervalServerError) {
				return nil, failure.New(logic.IntervalServerError)
			}
			return nil, failure.Wrap(err)
		}

		if err = u.variantCommand.Update(ctx, update.Variants()); err != nil {
			if errors.Is(err, service.IntervalServerError) {
				return nil, failure.New(logic.IntervalServerError)
			}
			return nil, failure.Wrap(err)
		}

		resp, err := u.uniqueQuery.Select(ctx, update.Unique().ID())
		if err != nil {
			if errors.Is(err, service.NotFound) {
				return nil, failure.New(logic.NotFound)
			}
			if errors.Is(err, service.IntervalServerError) {
				return nil, failure.New(logic.IntervalServerError)
			}
			return nil, failure.Wrap(err)
		}

		unique := resp.ToUniqueDinosaur()
		return &unique, nil
	})
}

func (u Unique) Delete(ctx context.Context, id model.UniqueDinosaurID) error {
	return logic.UseTransactioner0(ctx, func(ctx context.Context) error {
		if _, err := u.uniqueQuery.Select(ctx, id); err != nil {
			if errors.Is(err, service.NotFound); err != nil {
				return failure.New(logic.NotFound)
			}
			return failure.Wrap(err)
		}
		if err := u.uniqueCommand.Delete(ctx, id); err != nil {
			if errors.Is(err, service.IntervalServerError) {
				return failure.New(logic.IntervalServerError)
			}
			return failure.Wrap(err)
		}
		return nil
	})
}

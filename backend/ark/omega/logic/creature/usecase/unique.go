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
	variantCommand service.VariantsCommandRepository
}

func NewUnique(injector *do.Injector) (*Unique, error) {
	return &Unique{
		dinoCommand:    do.MustInvoke[service.DinosaurCommandRepository](injector),
		uniqueQuery:    do.MustInvoke[service.UniqueQueryRepository](injector),
		uniqueCommand:  do.MustInvoke[service.UniqueCommandRepository](injector),
		variantCommand: do.MustInvoke[service.VariantsCommandRepository](injector),
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
	uniques, err := u.uniqueQuery.List(ctx)
	if err != nil {
		if errors.Is(err, service.IntervalServerError) {
			return nil, failure.New(logic.IntervalServerError)
		}
		return nil, failure.Wrap(err)
	}
	return uniques, nil
}

func (u Unique) Create(ctx context.Context, create service.CreateCreature) (*model.UniqueDinosaur, error) {
	return logic.UseTransactioner(ctx, func(ctx context.Context) (*model.UniqueDinosaur, error) {
		d, err := u.dinoCommand.Insert(ctx, create.CreateDinosaur)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		unique, err := u.uniqueCommand.Insert(ctx, create.CreateUniqueDinosaur)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		v, err := u.variantCommand.Insert(ctx, create.CreateVariants)
		if err != nil {
			return nil, failure.Wrap(err)
		}

		dino := model.NewDinosaur(d.ID(), d.Name(), d.Health(), d.Melee())
		vs := lo.Map(v.Values(), func(item model.DinosaurVariant, _ int) model.DinosaurVariant {
			return model.NewDinosaurVariant(
				variantModel.NewVariant(item.ID(), item.Group(), item.Name()),
				model.VariantDescriptions{},
			)
		})

		resp := model.NewUniqueDinosaur(
			dino, unique.ID(), unique.Name(), vs, unique.HealthMultiplier(), unique.MeleeMultiplier(),
		)
		return &resp, nil
	})
}

func (u Unique) Update(ctx context.Context, update service.UpdateCreature) (*model.UniqueDinosaur, error) {
	return logic.UseTransactioner(ctx, func(ctx context.Context) (*model.UniqueDinosaur, error) {
		if _, err := u.uniqueQuery.Select(ctx, update.ID()); err != nil {
			if errors.Is(err, service.NotFound) {
				return nil, failure.New(logic.NotFound)
			}
			return nil, failure.Wrap(err)
		}

		d, err := u.dinoCommand.Update(ctx, update.UpdateDinosaur)
		if err != nil {
			if errors.Is(err, service.IntervalServerError) {
				return nil, failure.New(logic.IntervalServerError)
			}
			return nil, failure.Wrap(err)
		}

		unique, err := u.uniqueCommand.Update(ctx, update.UpdateUniqueDinosaur)
		if err != nil {
			if errors.Is(err, service.IntervalServerError) {
				return nil, failure.New(logic.IntervalServerError)
			}
			return nil, failure.Wrap(err)
		}

		v, err := u.variantCommand.Update(ctx, update.UpdateVariants)
		if err != nil {
			if errors.Is(err, service.IntervalServerError) {
				return nil, failure.New(logic.IntervalServerError)
			}
			return nil, failure.Wrap(err)
		}

		dino := model.NewDinosaur(d.ID(), d.Name(), d.Health(), d.Melee())
		vs := lo.Map(v.Values(), func(item model.DinosaurVariant, _ int) model.DinosaurVariant {
			return model.NewDinosaurVariant(
				variantModel.NewVariant(item.ID(), item.Group(), item.Name()),
				model.VariantDescriptions{},
			)
		})

		resp := model.NewUniqueDinosaur(
			dino, unique.ID(), unique.Name(), vs, unique.HealthMultiplier(), unique.MeleeMultiplier(),
		)

		return &resp, nil
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

package usecase

import (
	"context"
	"errors"

	"github.com/morikuni/failure"
	"github.com/samber/do"

	"mods-explore/ark/omega/logic"
	"mods-explore/ark/omega/logic/creature/domain/model"
	"mods-explore/ark/omega/logic/creature/domain/service"
)

type UniqueUsecase interface {
	Find(context.Context, model.UniqueDinosaurID) (*model.UniqueDinosaur, error)
	List(context.Context) (model.UniqueDinosaurs, error)
	Create(context.Context, service.CreateUniqueDinosaur) (*model.UniqueDinosaur, error)
	Update(context.Context, service.UpdateUniqueDinosaur) (*model.UniqueDinosaur, error)
	Delete(context.Context, model.UniqueDinosaurID) error
}

type Unique struct {
	repo service.UniqueRepository
}

func NewUnique(injector *do.Injector) (*Unique, error) {
	return &Unique{repo: do.MustInvoke[service.UniqueRepository](injector)}, nil
}

func (u Unique) Find(ctx context.Context, id model.UniqueDinosaurID) (*model.UniqueDinosaur, error) {
	unique, err := u.repo.Select(ctx, id)
	if err != nil {
		if errors.Is(err, service.NotFound) {
			return nil, failure.New(logic.NotFound)
		} else if errors.Is(err, service.IntervalServerError) {
			return nil, failure.New(logic.IntervalServerError)
		}
		return nil, failure.Wrap(err)
	}

	return unique, nil
}

func (u Unique) List(ctx context.Context) (model.UniqueDinosaurs, error) {
	uniques, err := u.repo.List(ctx)
	if err != nil {
		if errors.Is(err, service.IntervalServerError) {
			return nil, failure.New(logic.IntervalServerError)
		}
		return nil, failure.Wrap(err)
	}
	return uniques, nil
}

func (u Unique) Create(ctx context.Context, create service.CreateUniqueDinosaur) (*model.UniqueDinosaur, error) {
	return logic.UseTransactioner(ctx, func(ctx context.Context) (*model.UniqueDinosaur, error) {
		unique, err := u.repo.Insert(ctx, create)
		if err != nil {
			return nil, failure.Wrap(err)
		}
		return unique, nil
	})
}

func (u Unique) Update(ctx context.Context, update service.UpdateUniqueDinosaur) (*model.UniqueDinosaur, error) {
	return nil, nil
}

func (u Unique) Delete(ctx context.Context, id model.UniqueDinosaurID) error {
	return nil
}

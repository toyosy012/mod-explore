package usecase

import (
	"context"

	"github.com/samber/do"

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
	return nil, nil
}
func (u Unique) List(ctx context.Context) (model.UniqueDinosaurs, error) {
	return nil, nil
}
func (u Unique) Create(ctx context.Context, create service.CreateUniqueDinosaur) (*model.UniqueDinosaur, error) {
	return nil, nil
}
func (u Unique) Update(ctx context.Context, update service.UpdateUniqueDinosaur) (*model.UniqueDinosaur, error) {
	return nil, nil
}
func (u Unique) Delete(ctx context.Context, id model.UniqueDinosaurID) error {
	return nil
}

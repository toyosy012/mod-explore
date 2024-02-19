package storage

import (
	"context"

	"github.com/samber/do"

	"mods-explore/ark/omega/logic/creature/domain/model"
	"mods-explore/ark/omega/logic/creature/domain/service"
)

type DinosaurModel struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	BaseHealth int    `db:"health"`
	BaseMelee  int    `db:"melee"`
}

type DinosaurClient struct {
	*Client[DinosaurModel, int]
}

func NewDinosaurClient(injector *do.Injector) (service.DinosaurCommandRepository, error) {
	return DinosaurClient{
		do.MustInvoke[*Client[DinosaurModel, int]](injector),
	}, nil
}

func (c DinosaurClient) Insert(ctx context.Context, create service.CreateDinosaur) (*service.ResponseDinosaur, error) {
	return nil, nil
}

func (c DinosaurClient) Update(ctx context.Context, update service.UpdateDinosaur) (*service.ResponseDinosaur, error) {
	return nil, nil
}
func (c DinosaurClient) Delete(ctx context.Context, id model.DinosaurID) error {
	return nil
}

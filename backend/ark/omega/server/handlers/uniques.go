package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/samber/lo"

	creatureModel "mods-explore/ark/omega/logic/creature/domain/model"
	creatureSvc "mods-explore/ark/omega/logic/creature/domain/service"
	"mods-explore/ark/omega/logic/creature/usecase"
	variantModel "mods-explore/ark/omega/logic/variant/domain/model"
)

type UniqueHandler interface {
	ReadUnique(echo.Context) error
	ListUniques(echo.Context) error
	CreateUnique(echo.Context) error
	UpdateUnique(echo.Context) error
	DeleteUnique(echo.Context) error
}

type Unique struct {
	usecase.UniqueUsecase
}

func NewUnique(injector *do.Injector) (UniqueHandler, error) {
	return &Unique{
		UniqueUsecase: do.MustInvoke[usecase.UniqueUsecase](injector),
	}, nil
}

type uniqueQueryParams struct {
	ID int `param:"id" validate:"required"`
}

type UniqueValue struct {
	UniqueID         creatureModel.UniqueDinosaurID `json:"id" validate:"required"`
	BaseID           creatureModel.DinosaurID       `json:"base_id" validate:"required"`
	BaseName         creatureModel.DinosaurName     `json:"base_name" validate:"required"`
	BaseHealth       creatureModel.Health           `json:"base_health" validate:"required"`
	BaseMelee        creatureModel.Melee            `json:"base_melee" validate:"required"`
	UniqueName       creatureModel.UniqueName       `json:"unique_name" validate:"required"`
	HealthMultiplier float32                        `json:"health_multiplier" validate:"required"`
	DamageMultiplier float32                        `json:"damage_multiplier" validate:"required"`
	UniqueVariant    creatureModel.UniqueVariant    `json:"unique_variant" validate:"required"`
}

func NewUniqueValue(unique creatureModel.UniqueDinosaur) UniqueValue {
	return UniqueValue{
		unique.UniqueID(),
		unique.Dinosaur.BaseID(),
		unique.Dinosaur.BaseName(),
		unique.Dinosaur.Health(),
		unique.Dinosaur.Melee(),
		unique.UniqueName(),
		unique.HealthMultiplier().Value(),
		unique.DamageMultiplier().Value(),
		unique.UniqueVariant(),
	}
}

type UniqueValues []UniqueValue

func NewUniqueValues(uniques creatureModel.UniqueDinosaurs) UniqueValues {
	return lo.Map(uniques, func(u creatureModel.UniqueDinosaur, _ int) UniqueValue {
		return UniqueValue{}
	})
}

func (u Unique) ReadUnique(c echo.Context) error {
	var params uniqueQueryParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	unique, err := u.UniqueUsecase.Find(c.Request().Context(), creatureModel.UniqueDinosaurID(params.ID))
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, NewUniqueValue(*unique)); err != nil {
		return err
	}
	return nil
}

func (u Unique) ListUniques(c echo.Context) error {
	uniques, err := u.UniqueUsecase.List(c.Request().Context())
	if err != nil {
		return err
	}
	if err = c.JSON(http.StatusOK, NewUniqueValues(uniques)); err != nil {
		return err
	}
	return nil
}

type uniqueCreateParams struct {
	BaseName         creatureModel.DinosaurName                           `json:"base_name" validate:"required"`
	BaseHealth       creatureModel.Health                                 `json:"base_health" validate:"required"`
	BaseMelee        creatureModel.Melee                                  `json:"base_melee" validate:"required"`
	UniqueName       creatureModel.UniqueName                             `json:"unique_name" validate:"required"`
	HealthMultiplier creatureModel.UniqueMultiplier[creatureModel.Health] `json:"health_multiplier" validate:"required"`
	DamageMultiplier creatureModel.UniqueMultiplier[creatureModel.Melee]  `json:"damage_multiplier" validate:"required"`
	VariantIDs       [2]variantModel.VariantID                            `json:"variant_ids" validate:"required"`
}

func (u Unique) CreateUnique(c echo.Context) error {
	var params uniqueCreateParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	unique, err := u.UniqueUsecase.Create(
		c.Request().Context(),
		creatureSvc.NewCreateCreature(
			params.BaseName,
			params.BaseHealth,
			params.BaseMelee,
			params.UniqueName,
			params.HealthMultiplier,
			params.DamageMultiplier,
			params.VariantIDs,
		),
	)
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, NewUniqueValue(*unique)); err != nil {
		return err
	}
	return nil
}

type uniqueUpdateParams struct {
	UniqueID creatureModel.UniqueDinosaurID `param:"id" validate:"required"`

	BaseID           creatureModel.DinosaurID                             `json:"base_id" validate:"required"`
	BaseName         creatureModel.DinosaurName                           `json:"base_name" validate:"required"`
	BaseHealth       creatureModel.Health                                 `json:"base_health" validate:"required"`
	BaseMelee        creatureModel.Melee                                  `json:"base_melee" validate:"required"`
	UniqueName       creatureModel.UniqueName                             `json:"unique_name" validate:"required"`
	HealthMultiplier creatureModel.UniqueMultiplier[creatureModel.Health] `json:"health_multiplier" validate:"required"`
	DamageMultiplier creatureModel.UniqueMultiplier[creatureModel.Melee]  `json:"damage_multiplier" validate:"required"`
	UniqueVariantID  creatureModel.UniqueVariantID                        `json:"unique_variant_id" validate:"required"`
	VariantIDs       [2]variantModel.VariantID                            `json:"variant_ids" validate:"required"`
}

func (u Unique) UpdateUnique(c echo.Context) error {
	var params uniqueUpdateParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	unique, err := u.UniqueUsecase.Update(
		c.Request().Context(),
		creatureSvc.NewUpdateCreature(
			params.BaseID,
			params.BaseName,
			params.BaseHealth,
			params.BaseMelee,
			params.UniqueID,
			params.UniqueName,
			params.HealthMultiplier,
			params.DamageMultiplier,
			params.UniqueVariantID,
			params.VariantIDs,
		),
	)
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, NewUniqueValue(*unique)); err != nil {
		return err
	}
	return nil
}

func (u Unique) DeleteUnique(c echo.Context) error {
	var params variantGroupParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	err := u.UniqueUsecase.Delete(c.Request().Context(), creatureModel.UniqueDinosaurID(params.ID))
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, map[string]any{}); err != nil {
		return err
	}
	return nil
}

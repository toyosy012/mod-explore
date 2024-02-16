package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"

	"mods-explore/ark/omega/logic/variant/domain/model"
	"mods-explore/ark/omega/logic/variant/domain/service"
	"mods-explore/ark/omega/logic/variant/usecase"
)

type Variant struct {
	usecase.VariantUsecase
}

func NewVariant(injector *do.Injector) (VariantHandler, error) {
	return &Variant{
		VariantUsecase: do.MustInvoke[usecase.VariantUsecase](injector),
	}, nil
}

type VariantHandler interface {
	Read(echo.Context) error
	List(echo.Context) error
	Create(echo.Context) error
	Update(echo.Context) error
	Delete(echo.Context) error
}

type referenceParams struct {
	VariantID int `param:"id" validator:"required"`
}

func (v Variant) Read(c echo.Context) error {
	var params referenceParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	variant, err := v.VariantUsecase.Find(c.Request().Context(), model.VariantID(params.VariantID))
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, NewVariantValue(*variant)); err != nil {
		return err
	}
	return nil
}

func (v Variant) List(c echo.Context) error {
	variants, err := v.VariantUsecase.List(c.Request().Context())
	if err != nil {
		return err
	}
	if err = c.JSON(http.StatusOK, NewVariantValues(variants)); err != nil {
		return err
	}
	return nil
}

type createBody struct {
	GroupID int    `json:"group_id" validate:"required"`
	Name    string `json:"name" validate:"required"`
}

func (v Variant) Create(c echo.Context) error {
	var body createBody
	if err := c.Bind(&body); err != nil {
		return err
	}

	variant, err := v.VariantUsecase.Create(
		c.Request().Context(),
		service.NewCreateVariant(
			model.VariantGroupID(body.GroupID),
			model.Name(body.Name),
		),
	)
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, NewVariantValue(*variant)); err != nil {
		return err
	}
	return nil
}

type updateBody struct {
	VariantID int    `param:"id" validator:"required"`
	GroupID   int    `json:"group_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
}

func (v Variant) Update(c echo.Context) error {
	var body updateBody
	if err := c.Bind(&body); err != nil {
		return err
	}

	variant, err := v.VariantUsecase.Update(
		c.Request().Context(),
		service.NewUpdateVariant(model.VariantID(body.VariantID), model.VariantGroupID(body.GroupID), model.Name(body.Name)),
	)
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, NewVariantValue(*variant)); err != nil {
		return err
	}
	return nil
}

func (v Variant) Delete(c echo.Context) error {
	var params referenceParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	err := v.VariantUsecase.Delete(c.Request().Context(), model.VariantID(params.VariantID))
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, map[string]any{}); err != nil {
		return err
	}
	return nil
}

type VariantValue struct {
	ID    model.VariantID        `json:"id" validator:"required"`
	Name  model.Name             `json:"name" validator:"required"`
	Group model.VariantGroupName `json:"group" validator:"required"`
}

func NewVariantValue(v model.Variant) VariantValue {
	return VariantValue{
		ID:    v.ID(),
		Name:  v.Name(),
		Group: v.Group(),
	}
}

type VariantValues []VariantValue

func NewVariantValues(vs []model.Variant) (values VariantValues) {
	for _, v := range vs {
		values = append(
			values,
			VariantValue{
				ID:    v.ID(),
				Name:  v.Name(),
				Group: v.Group(),
			},
		)
	}
	return values
}

package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"mods-explore/ark/omega/logic/variant/domain/model"
	"mods-explore/ark/omega/logic/variant/domain/service"
	"mods-explore/ark/omega/logic/variant/usecase"
)

type VariantGroup struct {
	usecase.VariantGroupUsecase
}

func NewVariantGroup(usecase usecase.VariantGroupUsecase) VariantGroupHandler {
	return VariantGroup{
		usecase,
	}
}

type VariantGroupHandler interface {
	Read(echo.Context) error
	List(echo.Context) error
	Create(echo.Context) error
	Update(echo.Context) error
	Delete(echo.Context) error
}

type variantGroupParams struct {
	ID int `param:"id" validator:"required"`
}

func (v VariantGroup) Read(c echo.Context) error {
	var params variantGroupParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	variantGroup, err := v.VariantGroupUsecase.Find(c.Request().Context(), model.VariantGroupID(params.ID))
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, NewVariantGroupValue(*variantGroup)); err != nil {
		return err
	}
	return nil
}

func (v VariantGroup) List(c echo.Context) error {
	variantGroups, err := v.VariantGroupUsecase.List(c.Request().Context())
	if err != nil {
		return err
	}
	if err = c.JSON(http.StatusOK, NewVariantGroupValues(variantGroups)); err != nil {
		return err
	}
	return nil
}

type createVariantGroup struct {
	Name string `json:"name" validate:"required"`
}

func (v VariantGroup) Create(c echo.Context) error {
	var body createVariantGroup
	if err := c.Bind(&body); err != nil {
		return err
	}

	variant, err := v.VariantGroupUsecase.Create(
		c.Request().Context(),
		service.NewCreateVariantGroup(model.VariantGroupName(body.Name)),
	)
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, NewVariantGroupValue(*variant)); err != nil {
		return err
	}
	return nil
}

type updateVariantGroup struct {
	ID   int    `param:"id" validator:"required"`
	Name string `json:"name" validate:"required"`
}

func (v VariantGroup) Update(c echo.Context) error {
	var body updateVariantGroup
	if err := c.Bind(&body); err != nil {
		return err
	}

	variant, err := v.VariantGroupUsecase.Update(
		c.Request().Context(),
		service.NewUpdateVariantGroup(model.VariantGroupID(body.ID), model.VariantGroupName(body.Name)),
	)
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, NewVariantGroupValue(*variant)); err != nil {
		return err
	}
	return nil
}

func (v VariantGroup) Delete(c echo.Context) error {
	var params variantGroupParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	err := v.VariantGroupUsecase.Delete(c.Request().Context(), model.VariantGroupID(params.ID))
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, map[string]any{}); err != nil {
		return err
	}
	return nil
}

type VariantGroupValue struct {
	ID   model.VariantGroupID   `json:"id" validator:"required"`
	Name model.VariantGroupName `json:"name" validator:"required"`
}

func NewVariantGroupValue(v model.VariantGroup) VariantGroupValue {
	return VariantGroupValue{
		ID:   v.ID(),
		Name: v.Name(),
	}
}

type VariantGroupValues []VariantGroupValue

func NewVariantGroupValues(vs []model.VariantGroup) (values VariantGroupValues) {
	for _, v := range vs {
		values = append(
			values,
			VariantGroupValue{
				ID:   v.ID(),
				Name: v.Name(),
			},
		)
	}
	return values
}

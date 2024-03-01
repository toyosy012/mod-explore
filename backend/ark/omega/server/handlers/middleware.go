package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/morikuni/failure"
	"github.com/samber/do"

	"mods-explore/ark/omega/logic"
	"mods-explore/ark/omega/storage"
)

func NewErrorHandler(s *echo.Echo) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		code, ok := failure.CodeOf(err)
		if !ok {
			s.DefaultHTTPErrorHandler(err, c)
			return
		}

		switch code {
		case logic.InvalidArgument:
			if err := c.JSON(http.StatusBadRequest, map[string]any{
				"message": "bad request",
			}); err != nil {
				c.Logger().Error(err)
			}
			return
		case logic.NotFound:
			if err := c.JSON(http.StatusNotFound, map[string]any{
				"message": "not found",
			}); err != nil {
				c.Logger().Error(err)
			}
			return
		case logic.Forbidden:
			if err := c.NoContent(http.StatusForbidden); err != nil {
				if err != nil {
					c.Logger().Error(err)
				}
			}
		default:
			if err := c.NoContent(http.StatusInternalServerError); err != nil {
				if err != nil {
					c.Logger().Error(err)
				}
			}
		}
	}
}

func Transctioner(injector *do.Injector) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = logic.SetTransactioner(ctx, do.MustInvoke[*storage.Client](injector))
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

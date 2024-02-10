package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/morikuni/failure"

	"mods-explore/ark/omega/logic"
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
		status = http.StatusBadRequest
	case logic.NotFound:
		status = http.StatusNotFound
	case logic.Forbidden:
		status = http.StatusForbidden
	default:
		status = http.StatusInternalServerError
	}

	if err = c.JSON(status, err.Error()); err != nil {
		c.Logger().Error(err)
	}
}

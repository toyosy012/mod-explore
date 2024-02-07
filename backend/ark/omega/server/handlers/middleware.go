package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/morikuni/failure"

	"mods-explore/ark/omega/logic"
)

func ErrorHandler(err error, c echo.Context) {
	var status int
	code, ok := failure.CodeOf(err)
	if !ok {
		status = http.StatusInternalServerError
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

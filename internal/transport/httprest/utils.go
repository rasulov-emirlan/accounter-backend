package httprest

import (
	"github.com/labstack/echo/v4"
)

func respondErr(ctx echo.Context, code int, err error) error {
	return ctx.JSON(code, echo.Map{
		"error": err.Error(),
	})
}

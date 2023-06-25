package httprest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/esep-backend/pkg/validation"
)

func respondErr(ctx echo.Context, code int, err error) error {
	if code == http.StatusBadRequest {
		return ctx.JSON(code, echo.Map{
			"error": validation.GetValidator().Mappify(err),
		})
	}
	return ctx.JSON(code, echo.Map{
		"error": err.Error(),
	})
}

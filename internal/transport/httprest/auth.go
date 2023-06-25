package httprest

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/esep-backend/internal/domains/auth"
	"github.com/rasulov-emirlan/esep-backend/pkg/validation"
)

const AuthRefreshCookieName = "refresh_token"

type AuthRefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type AuthHandler struct {
	service auth.Service
}

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

func (h AuthHandler) Register(ctx echo.Context) error {
	req := new(auth.RegisterInput)
	if err := ctx.Bind(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	if err := ctx.Validate(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	session, err := h.service.Register(ctx.Request().Context(), *req)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	ctx.SetCookie(&http.Cookie{
		Name:     AuthRefreshCookieName,
		Value:    session.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return ctx.JSON(http.StatusOK, session)
}

func (h AuthHandler) Login(ctx echo.Context) error {
	req := new(auth.LoginInput)
	if err := ctx.Bind(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	if err := ctx.Validate(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	session, err := h.service.Login(ctx.Request().Context(), *req)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	ctx.SetCookie(&http.Cookie{
		Name:     AuthRefreshCookieName,
		Value:    session.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return ctx.JSON(http.StatusOK, session)
}

func (h AuthHandler) Refresh(ctx echo.Context) error {
	refreshToken, err := ctx.Cookie(AuthRefreshCookieName)
	if err != nil {
		// search for refresh token in request body
		req := new(AuthRefreshRequest)
		if err := ctx.Bind(req); err != nil {
			return respondErr(ctx, http.StatusBadRequest, err)
		}
		if ctx.Validate(req) == nil {
			refreshToken = &http.Cookie{
				Value: req.RefreshToken,
			}
		}
		return respondErr(ctx, http.StatusBadRequest, errors.New("refresh token is required in cookie or request body"))
	}

	session, err := h.service.Refresh(ctx.Request().Context(), refreshToken.Value)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	ctx.SetCookie(&http.Cookie{
		Name:     AuthRefreshCookieName,
		Value:    session.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return ctx.JSON(http.StatusOK, session)
}

func (h AuthHandler) Logout(ctx echo.Context) error {

	ctx.SetCookie(&http.Cookie{
		Name:     AuthRefreshCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return ctx.NoContent(http.StatusOK)
}

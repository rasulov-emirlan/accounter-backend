package httprest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/esep-backend/internal/domains/auth"
)

const AuthRefreshCookieName = "refresh_token"

type AuthHandler struct {
	service auth.Service
}

func respondErr(ctx echo.Context, code int, err error) error {
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
		return respondErr(ctx, http.StatusBadRequest, err)
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

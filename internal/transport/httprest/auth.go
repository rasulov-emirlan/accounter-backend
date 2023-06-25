package httprest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/esep-backend/internal/domains/auth"
)

const (
	AuthRefreshCookieName  = "refresh_token"
	AuthSessionContextName = "session"
)

type AuthRefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type AuthHandler struct {
	service auth.Service
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

func (h AuthHandler) MiddlewareUnpackAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		accessHeader := ctx.Request().Header.Get("Authorization")
		if accessHeader == "" {
			return respondErr(ctx, http.StatusUnauthorized, errors.New("access token is required in 'Authorization' header"))
		}

		access := strings.Split(accessHeader, " ")
		if len(access) != 2 || access[0] != "Bearer" {
			return respondErr(ctx, http.StatusUnauthorized, errors.New("invalid access token format"))
		}

		session, err := h.service.ParseAccessKey(ctx.Request().Context(), access[1])
		if err != nil {
			return respondErr(ctx, http.StatusUnauthorized, err)
		}

		ctx.Set(AuthSessionContextName, session)
		return next(ctx)
	}
}

func (h AuthHandler) Me(ctx echo.Context) error {
	session, ok := ctx.Get(AuthSessionContextName).(auth.AccessKey)
	if !ok {
		return respondErr(ctx, http.StatusInternalServerError, errors.New("session not found in context"))
	}

	me, err := h.service.Me(ctx.Request().Context(), session)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, me)
}

package httprest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rasulov-emirlan/esep-backend/config"
	"github.com/rasulov-emirlan/esep-backend/internal/domains"
	"github.com/rasulov-emirlan/esep-backend/pkg/health"
	"github.com/rasulov-emirlan/esep-backend/pkg/logging"
)

var (
	ErrServerClosed = http.ErrServerClosed
)

type server struct {
	srvr *http.Server
}

func NewServer(cfg config.Config) *server {
	return &server{
		srvr: &http.Server{
			Addr:         cfg.Server.Port,
			ReadTimeout:  cfg.Server.TimeoutRead,
			WriteTimeout: cfg.Server.TimeoutWrite,
		},
	}
}

func (s *server) Start(log *logging.Logger, doms domains.DomainCombiner) error {
	router := echo.New()

	router.Use(middleware.Recover())
	router.Use(middleware.CORS())
	router.Use(middleware.Secure())
	router.Use(middleware.AddTrailingSlash())
	router.Use(middleware.Gzip())
	router.Use(log.NewEchoMiddleware)

	router.Any(
		"/health*",
		echo.WrapHandler(health.NewHTTPHandler("esep-backend", nil)),
	)

	authHandler := AuthHandler{doms.AuthService()}
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
	}

	s.srvr.Handler = router

	return s.srvr.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	return s.srvr.Shutdown(ctx)
}

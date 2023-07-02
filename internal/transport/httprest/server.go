package httprest

import (
	"context"
	"net/http"

	"github.com/SanaripEsep/esep-backend/config"
	"github.com/SanaripEsep/esep-backend/internal/domains"
	"github.com/SanaripEsep/esep-backend/pkg/health"
	"github.com/SanaripEsep/esep-backend/pkg/logging"
	"github.com/SanaripEsep/esep-backend/pkg/validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	router.Use(middleware.RemoveTrailingSlash())
	router.Use(middleware.Gzip())
	router.Use(log.NewEchoMiddleware)
	router.Validator = validation.GetValidator()

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
		authGroup.GET("/me", authHandler.Me, authHandler.MiddlewareUnpackAccess)
	}

	storesHandler := StoresHandler{doms.StoresService()}
	storesGroup := router.Group("/stores", authHandler.MiddlewareUnpackAccess)
	{
		storesGroup.GET("/:id", storesHandler.Read)
		storesGroup.GET("", storesHandler.ReadBy)
		storesGroup.POST("", storesHandler.Create)
		storesGroup.PATCH("/:id", storesHandler.Update)
		storesGroup.DELETE("/:id", storesHandler.Delete)
	}

	categoriesHandler := CategoriesHandler{doms.CategoriesService()}
	categoriesGroup := router.Group("/categories", authHandler.MiddlewareUnpackAccess)
	{
		categoriesGroup.GET("/:id", categoriesHandler.Read)
		categoriesGroup.GET("", categoriesHandler.ReadBy)
		categoriesGroup.POST("", categoriesHandler.Create)
		categoriesGroup.PATCH("/:id", categoriesHandler.Update)
		categoriesGroup.DELETE("/:id", categoriesHandler.Delete)
	}

	s.srvr.Handler = router

	return s.srvr.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	return s.srvr.Shutdown(ctx)
}

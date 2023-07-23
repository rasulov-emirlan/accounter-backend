package httprest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"github.com/rasulov-emirlan/accounter-backend/config"
	"github.com/rasulov-emirlan/accounter-backend/internal/domains"
	"github.com/rasulov-emirlan/accounter-backend/pkg/health"
	"github.com/rasulov-emirlan/accounter-backend/pkg/logging"
	"github.com/rasulov-emirlan/accounter-backend/pkg/telemetry"
	"github.com/rasulov-emirlan/accounter-backend/pkg/validation"
)

var (
	ErrServerClosed = http.ErrServerClosed
)

type server struct {
	srvr        *http.Server
	serviceName string
}

func NewServer(cfg config.Config) server {
	return server{
		srvr: &http.Server{
			Addr:         cfg.Server.Port,
			ReadTimeout:  cfg.Server.TimeoutRead,
			WriteTimeout: cfg.Server.TimeoutWrite,
		},
		serviceName: cfg.ServiceName,
	}
}

func (s server) Start(log *logging.Logger, doms domains.DomainCombiner) error {
	router := echo.New()

	router.Use(middleware.Recover())
	router.Use(middleware.CORS())
	router.Use(middleware.Secure())
	router.Use(middleware.RemoveTrailingSlash())
	router.Use(middleware.Gzip())
	router.Use(log.NewEchoMiddleware)
	router.Use(otelecho.Middleware(s.serviceName))
	router.HTTPErrorHandler = telemetry.EchoHTTPErrorHandler(router)
	router.Validator = validation.GetValidator()

	router.Any(
		"/health*",
		echo.WrapHandler(health.NewHTTPHandler(s.serviceName, nil)),
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

func (s server) Stop(ctx context.Context) error {
	return s.srvr.Shutdown(ctx)
}

package httprest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/esep-backend/pkg/health"
	"github.com/rasulov-emirlan/esep-backend/pkg/logging"
)

var (
	ErrServerClosed = http.ErrServerClosed
)

type server struct {
	srvr *http.Server
}

func NewServer(addr string) *server {
	return &server{
		srvr: &http.Server{
			Addr: addr,
		},
	}
}

func (s *server) Start(log *logging.Logger) error {
	router := echo.New()

	router.Use(log.NewEchoMiddleware)

	router.Any(
		"/health*",
		echo.WrapHandler(health.NewHTTPHandler("esep-backend", nil)),
	)

	s.srvr.Handler = router

	return s.srvr.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	return s.srvr.Shutdown(ctx)
}

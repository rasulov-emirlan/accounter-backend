package telemetry

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/SanaripEsep/esep-backend/config"
	"github.com/SanaripEsep/esep-backend/pkg/shutdown"
)

const SpanNamePrefix = "github.com/SanaripEsep/esep-backend/"

func Name(name string) string {
	return SpanNamePrefix + name
}

func NewSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(name).Start(ctx, name)
	return span
}

func StartJaegerTraceProvider(cfg config.Config, cleaner *shutdown.Scheduler) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JeagerURL)))
	if err != nil {
		return err
	}

	envKey := "production"
	if cfg.Flags.DevMode {
		envKey = "development"
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.DeploymentEnvironmentKey.String(envKey),
		)),
	)

	cleaner.Add(tp.Shutdown)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}

func EchoHTTPErrorHandler(router *echo.Echo) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		ctx := c.Request().Context()
		trace.SpanFromContext(ctx).RecordError(err)

		router.DefaultHTTPErrorHandler(err, c)
	}
}

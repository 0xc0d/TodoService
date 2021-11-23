package trace

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"time"
)

type Shutdown func(context.Context)

// tracerProvider returns an OpenTelemetry TracerProvider configured to
// use the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with given service name and
// environment.
func tracerProvider(url, service, environment string) (*trace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
		)),
	)

	return tp, nil
}

func SetupJaeger(endpoint, service, env string) Shutdown {
	tp, err := tracerProvider(endpoint, service, env)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	otel.SetTracerProvider(tp)

	// Cleanly shutdown and flush telemetry when the application exits.
	return func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err = tp.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Send()
		}
	}
}

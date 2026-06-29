package config

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var TP *sdktrace.TracerProvider

func InitTracing() {
	ctx := context.Background()

	var opts []otlptracehttp.Option
	if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
		opts = append(opts, otlptracehttp.WithEndpoint(endpoint))
	}
	if insecure := os.Getenv("OTEL_EXPORTER_OTLP_INSECURE"); insecure == "true" {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		panic("PANIC :: failed to create OTLP trace exporter: " + err.Error())
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "mutt"
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		panic("PANIC :: failed to create OTel resource: " + err.Error())
	}

	TP = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(TP)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

// TODO : Need to check if graceful shutdown is implemented for the tracer provider.
// If not, implement it.
func ShutdownTracing() {
	if TP != nil {
		TP.Shutdown(context.Background())
	}
}

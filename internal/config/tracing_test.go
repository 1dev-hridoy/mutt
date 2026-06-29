package config

import (
	"context"
	"os"
	"testing"

	"go.opentelemetry.io/otel"
)

func TestInitTracing_SetsGlobalProvider(t *testing.T) {
	os.Setenv("OTEL_SERVICE_NAME", "mutt-test")
	os.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "true")
	defer os.Unsetenv("OTEL_SERVICE_NAME")
	defer os.Unsetenv("OTEL_EXPORTER_OTLP_INSECURE")

	providerBefore := otel.GetTracerProvider()

	InitTracing()

	providerAfter := otel.GetTracerProvider()
	if providerAfter == providerBefore {
		t.Fatal("expected global TracerProvider to be updated after InitTracing")
	}
	if TP == nil {
		t.Fatal("expected TP to be set after InitTracing")
	}
}

func TestShutdownTracing_NoPanicWhenNil(t *testing.T) {
	TP = nil
	ShutdownTracing()
}

func TestInitTracing_DefaultServiceName(t *testing.T) {
	os.Unsetenv("OTEL_SERVICE_NAME")
	os.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "true")
	defer os.Unsetenv("OTEL_EXPORTER_OTLP_INSECURE")

	InitTracing()

	if TP == nil {
		t.Fatal("expected TP to be set")
	}
	TP.Shutdown(context.Background())
}

func TestInitTracing_CustomEndpoint(t *testing.T) {
	os.Setenv("OTEL_SERVICE_NAME", "mutt-test")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:9999")
	os.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "true")
	defer os.Unsetenv("OTEL_SERVICE_NAME")
	defer os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	defer os.Unsetenv("OTEL_EXPORTER_OTLP_INSECURE")

	InitTracing()

	if TP == nil {
		t.Fatal("expected TP to be set")
	}
	TP.Shutdown(context.Background())
}

package config

import (
	"io"
	"os"

	gcpexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
)

func initTracerExporter(cfg *Config) (sdktrace.SpanExporter, error) {
	switch cfg.Trace.Exporter {
	case "jaeger":
		// Create the Jaeger exporter
		return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Trace.Jaeger.Endpoint)))
	case "gcp":
		projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
		return gcpexporter.New(gcpexporter.WithProjectID(projectID))
	case "stdout":
		return stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithWriter(os.Stderr),
		)
	case "none":
		return stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithWriter(io.Discard),
		)
	default:
		return nil, libD.ErrInvalidArgument
	}
}

func InitTracerProvider(cfg *Config) (*sdktrace.TracerProvider, error) {
	exp, err := initTracerExporter(cfg)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.App.Name),
		)),
	)

	return tp, nil
}

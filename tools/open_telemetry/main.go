package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kujilabo/cocotola-api/pkg_app/config"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/gateway"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// var tracer = otel.Tracer("tracer")

const (
	service     = "trace-demo"
	environment = "production"
	id          = 1
)

func main() {
	ctx := context.Background()

	// tp := initTracer()
	// defer func() {
	// 	if err := tp.Shutdown(context.Background()); err != nil {
	// 		log.Printf("Error shutting down tracer provider: %v", err)
	// 	}
	// }()

	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	cfg, err := config.LoadConfig("local")
	if err != nil {
		panic(err)
	}

	client := gateway.NewTatoebaClient(cfg.Tatoeba.Endpoint, cfg.Tatoeba.Username, cfg.Tatoeba.Password, time.Duration(cfg.Tatoeba.Timeout)*time.Second)

	tr := tp.Tracer("component-main")
	_, span := tr.Start(ctx, "find", oteltrace.WithAttributes(attribute.Int("sentencenumber", 1286)))
	defer span.End()

	sentence, err := client.FindSentenceBySentenceNumber(ctx, 1296)
	if err != nil {
		panic(err)
	}

	fmt.Println(sentence.GetText())
	bar(ctx)
}

func bar(ctx context.Context) {
	// Use the global TracerProvider.
	tr := otel.Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()

	// Do bar...
}
func tracerProvider(url string) (*sdktrace.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

// func initTracer() *sdktrace.TracerProvider {
// 	exporter, err := stdout.New(stdout.WithPrettyPrint())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tp := sdktrace.NewTracerProvider(
// 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 		sdktrace.WithBatcher(exporter),
// 	)
// 	otel.SetTracerProvider(tp)
// 	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
// 	return tp
// }

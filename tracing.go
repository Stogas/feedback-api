package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Stogas/feedback-api/internal/config"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer(conf config.TraceConfig) (func(), propagation.TextMapPropagator) {
	ctx := context.Background()

	// GRPC Exporter
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%v", conf.Host, conf.Port)),
	)
	if err != nil {
		slog.Error("failed to initialize tracer", "error", err)
		panic("failed to initialize tracer")
	}

	// HTTP Exporter
	// exporter, err := otlptracehttp.New(ctx, otlp)
	// if err != nil {
	//     return nil, err
	// }

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("FeedbackAPI"),
			semconv.TelemetrySDKLanguageGo,
		)),
	)

	b3Propagator := b3.New()
	otel.SetTextMapPropagator(b3Propagator)

	otel.SetTracerProvider(tp)

	return func() {
		slog.Debug("Shutting down OTLP trace provider ...")
		if err := tp.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shut down tracer", "error", err)
			panic("failed to shut down tracer")
		}
		slog.Debug("OTLP trace provider exited successfully")
	}, b3Propagator
}

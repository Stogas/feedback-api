package main

import (
	"context"
	"fmt"

	"github.com/Stogas/feedback-api/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer(conf config.TraceConfig) (*trace.TracerProvider, error) {
	ctx := context.Background()

	// GRPC Exporter
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%v", conf.Host, conf.Port)),
	)
	if err != nil {
		return nil, err
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

	otel.SetTracerProvider(tp)
	return tp, nil
}

package telemetry

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type ApplicationMetrics struct {
	RunCounter     metric.Int64Counter
	GoRunsCounter  metric.Int64Counter
	JSRunsCounter  metric.Int64Counter
	SuccessCounter metric.Int64Counter
	FailureCounter metric.Int64Counter
	RunDuration    metric.Float64Histogram

	// Add more metrics as needed

}

func NewApplicationMetrics(meter metric.Meter) *ApplicationMetrics {

	runCounter, err := meter.Int64Counter("run_total")
	if err != nil {
		panic(fmt.Sprintf("failed to create run counter: %v", err))
	}

	goRunsCounter, err := meter.Int64Counter("run_go_total")
	if err != nil {
		panic(fmt.Sprintf("failed to create Go run counter: %v", err))
	}

	jsRunsCounter, err := meter.Int64Counter("run_js_total")
	if err != nil {
		panic(fmt.Sprintf("failed to create JS run counter: %v", err))
	}

	successCounter, err := meter.Int64Counter("run_success_total")
	if err != nil {
		panic(fmt.Sprintf("failed to create success counter: %v", err))
	}

	failureCounter, err := meter.Int64Counter("run_failure_total")
	if err != nil {
		panic(fmt.Sprintf("failed to create failure counter: %v", err))
	}

	runDuration, err := meter.Float64Histogram("run_duration_seconds")
	if err != nil {
		panic(fmt.Sprintf("failed to create run duration histogram: %v", err))
	}

	return &ApplicationMetrics{
		RunCounter:     runCounter,
		GoRunsCounter:  goRunsCounter,
		JSRunsCounter:  jsRunsCounter,
		SuccessCounter: successCounter,
		FailureCounter: failureCounter,
		RunDuration:    runDuration,
	}
}

func InitMetrics(ctx context.Context, collectorEndpoint string) (metric.Meter, func(context.Context) error, error) {
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithEndpoint(collectorEndpoint), otlpmetricgrpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
	)

	meter := mp.Meter("lgtm/backend")
	return meter, mp.Shutdown, nil
}

func InitTracing(ctx context.Context, collectorEndpoint string) (trace.Tracer, func(context.Context) error, error) {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(collectorEndpoint), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("lgtm:backend"),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	tracer := tp.Tracer("lgtm/backend")
	return tracer, tp.Shutdown, nil
}

func InitLogging(ctx context.Context, collectorEndpoint string) (*slog.Logger, func(context.Context) error, error) {
	exporter, err := otlploggrpc.New(ctx, otlploggrpc.WithEndpoint(collectorEndpoint), otlploggrpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("lgtm:backend"),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)

	handler := otelslog.NewHandler("lgtm/backend", otelslog.WithLoggerProvider(lp))
	logger := slog.New(handler)
	return logger, lp.Shutdown, nil
}

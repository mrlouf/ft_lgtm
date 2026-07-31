package backend

import (
	"context"
	"log"
	"time"

	"lgtm/internal/publisher"
	"lgtm/internal/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Compiler interface {
	Compile(ctx context.Context, source []byte, lang string) ([]byte, error)
}

type Executor interface {
	Execute(ctx context.Context, wasmBinary []byte) (stdout, stderr string, err error)
}

type Publisher interface {
	Publish(ctx context.Context, source []byte, stdout []byte) (response publisher.ResponseCID, err error)
}

type Backend struct {
	Compiler  Compiler
	Executor  Executor
	Publisher Publisher

	// TODO: Add Tracer/Metrics/Logger as an Obs struct
	Tracer  trace.Tracer
	Metrics *telemetry.ApplicationMetrics
	// Logger  *log.Logger
}

func NewBackend(
	compiler Compiler,
	executor Executor,
	publisher Publisher,
	tracer trace.Tracer,
	meter metric.Meter) *Backend {

	metrics := telemetry.NewApplicationMetrics(meter)

	return &Backend{
		Compiler:  compiler,
		Executor:  executor,
		Publisher: publisher,
		Tracer:    tracer,
		Metrics:   metrics,
	}
}

func (b *Backend) Run(ctx context.Context, source []byte, language string) (string, string, publisher.ResponseCID, error) {

	start := time.Now()
	ctx, span := b.Tracer.Start(ctx, "backend.run")
	defer span.End()

	log.Printf("Run: start for language: %s", language)
	b.Metrics.RunCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("language", language)))

	wasmBinary, err := b.Compiler.Compile(ctx, source, language)
	if err != nil {
		span.RecordError(err)
		b.Metrics.FailureCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("language", language), attribute.String("error_type", "compile")))
		b.Metrics.RunDuration.Record(ctx, time.Since(start).Seconds())
		return "", "", publisher.ResponseCID{}, err
	}

	stdout, stderr, err := b.Executor.Execute(ctx, wasmBinary)
	if err != nil {
		span.RecordError(err)
		b.Metrics.FailureCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("language", language), attribute.String("error_type", "execute")))
		b.Metrics.RunDuration.Record(ctx, time.Since(start).Seconds())
		return stdout, stderr, publisher.ResponseCID{}, err
	}

	responseCID, err := b.Publisher.Publish(ctx, source, []byte(stdout))
	if err != nil {
		span.RecordError(err)
		b.Metrics.FailureCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("language", language), attribute.String("error_type", "publish")))
		b.Metrics.RunDuration.Record(ctx, time.Since(start).Seconds())
		return stdout, stderr, publisher.ResponseCID{}, err
	}

	status := classifyStatus(err)

	log.Println("Run: completed with status:", status.String())

	b.Metrics.SuccessCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("language", language), attribute.String("status", status.String())))
	b.Metrics.RunDuration.Record(ctx, time.Since(start).Seconds())

	return stdout, stderr, responseCID, nil

}

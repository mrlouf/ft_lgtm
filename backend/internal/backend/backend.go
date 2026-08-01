package backend

import (
	"context"
	"log"
	"log/slog"
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
	Logger  *slog.Logger
}

func NewBackend(

	compiler Compiler,
	executor Executor,
	publisher Publisher,
	tracer trace.Tracer,
	meter metric.Meter,
	logger *slog.Logger,

) *Backend {

	metrics := telemetry.NewApplicationMetrics(meter)

	return &Backend{
		Compiler:  compiler,
		Executor:  executor,
		Publisher: publisher,
		Tracer:    tracer,
		Metrics:   metrics,
		Logger:    logger,
	}
}

type RunSpecs struct {
	Language string
	Source   []byte
	Start    time.Time
	Ctx      context.Context
	Span     trace.Span
}

func (b *Backend) Run(r RunSpecs) (string, string, publisher.ResponseCID, error) {

	log.Printf("Run: start for language: %s", r.Language)
	b.Metrics.RunCounter.Add(r.Ctx, 1, metric.WithAttributes(attribute.String("language", r.Language)))

	wasmBinary, err := b.Compiler.Compile(r.Ctx, r.Source, r.Language)
	if err != nil {
		r.Span.RecordError(err)
		b.Metrics.FailureCounter.Add(r.Ctx, 1, metric.WithAttributes(attribute.String("language", r.Language), attribute.String("error_type", "compile")))
		b.Metrics.RunDuration.Record(r.Ctx, time.Since(r.Start).Seconds())
		return "", "", publisher.ResponseCID{}, err
	}

	stdout, stderr, err := b.Executor.Execute(r.Ctx, wasmBinary)
	if err != nil {
		r.Span.RecordError(err)
		b.Metrics.FailureCounter.Add(r.Ctx, 1, metric.WithAttributes(attribute.String("language", r.Language), attribute.String("error_type", "execute")))
		b.Metrics.RunDuration.Record(r.Ctx, time.Since(r.Start).Seconds())
		return stdout, stderr, publisher.ResponseCID{}, err
	}

	responseCID, err := b.Publisher.Publish(r.Ctx, r.Source, []byte(stdout))
	if err != nil {

		b.Metrics.FailureCounter.Add(r.Ctx, 1, metric.WithAttributes(attribute.String("language", r.Language), attribute.String("error_type", "publish")))
		b.Metrics.RunDuration.Record(r.Ctx, time.Since(r.Start).Seconds())
		return stdout, stderr, publisher.ResponseCID{}, err
	}

	status := classifyStatus(err)

	log.Println("Run: completed with status:", status.String())

	b.Metrics.SuccessCounter.Add(r.Ctx, 1, metric.WithAttributes(attribute.String("language", r.Language), attribute.String("status", status.String())))
	b.Metrics.RunDuration.Record(r.Ctx, time.Since(r.Start).Seconds())

	return stdout, stderr, responseCID, nil

}

package backend

import (
	"context"
	"log/slog"
	"time"

	"lgtm/internal/publisher"
	"lgtm/internal/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Compiler interface {
	Compile(ctx context.Context, span trace.Span, source []byte, lang string) ([]byte, error)
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
	Span     trace.Span
}

func (b *Backend) Run(ctx context.Context, r RunSpecs) (string, string, publisher.ResponseCID, error) {

	ctx, span := b.Tracer.Start(ctx, "backend.run", trace.WithAttributes(
		attribute.String("language", r.Language),
	))
	defer span.End()
	b.Logger.InfoContext(ctx, "run: start", "language", r.Language)
	b.Metrics.RunCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("language", r.Language)))

	wasmBinary, err := b.Compiler.Compile(ctx, span, r.Source, r.Language)
	if err != nil {
		span.RecordError(err)
		b.Metrics.FailureCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("language", r.Language), attribute.String("error_type", "compile")))
		b.Metrics.RunDuration.Record(ctx, time.Since(r.Start).Seconds())
		return "", "", publisher.ResponseCID{}, err
	}

	stdout, stderr, err := b.Executor.Execute(ctx, wasmBinary)
	if err != nil {
		span.RecordError(err)
		b.Metrics.FailureCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("language", r.Language), attribute.String("error_type", "execute")))
		b.Metrics.RunDuration.Record(ctx, time.Since(r.Start).Seconds())
		return stdout, stderr, publisher.ResponseCID{}, err
	}

	responseCID, err := b.Publisher.Publish(ctx, r.Source, []byte(stdout))
	if err != nil {
		span.RecordError(err)
		b.Metrics.FailureCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("language", r.Language), attribute.String("error_type", "publish")))
		b.Metrics.RunDuration.Record(ctx, time.Since(r.Start).Seconds())
		return stdout, stderr, publisher.ResponseCID{}, err
	}

	status := classifyStatus(err)

	b.Logger.InfoContext(ctx, "run: completed", "language", r.Language, "status", status.String())
	span.AddEvent("run completed", trace.WithAttributes(
		attribute.String("status", status.String()),
	))
	b.Metrics.SuccessCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("language", r.Language), attribute.String("status", status.String())))
	b.Metrics.RunDuration.Record(ctx, time.Since(r.Start).Seconds())

	return stdout, stderr, responseCID, nil

}

package api

import (
	"context"
	"encoding/json"
	"lgtm/internal/backend"
	"lgtm/internal/publisher"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Request struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type Response struct {
	SourceCID string `json:"source_cid"`
	OutputCID string `json:"output_cid"`
	Status    string `json:"status"`
	Stdout    string `json:"stdout,omitempty"`
	Stderr    string `json:"stderr,omitempty"`
	Error     string `json:"error,omitempty"`
}

type Metrics struct {
	TotalRequests prometheus.Counter `json:"total_requests"`
	Successful    prometheus.Counter `json:"successful"`
	Failed        prometheus.Counter `json:"failed"`
}

func getHTTPStatusFromError(err error) int {

	log.Printf("Classifying error: %v", err)

	stage, _, _ := strings.Cut(err.Error(), ": ")

	switch stage {
	case "compile":
		return http.StatusBadRequest
	case "execute":
		return http.StatusBadRequest
	case "timeout":
		return http.StatusRequestTimeout
	case "publish":
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func returnFailedResponse(logger *slog.Logger, tracer trace.Tracer, w http.ResponseWriter, stderr string, err error) {

	ctx := context.Background()
	_, span := tracer.Start(ctx, "backend.returnFailedResponse")
	defer span.End()

	logger.Error("Run failed", slog.Any("error", err))

	w.WriteHeader(getHTTPStatusFromError(err))

	resp := Response{
		Status: "failed",
		Stderr: stderr,
		Error:  err.Error(),
	}

	httpStatus := getHTTPStatusFromError(err)
	span.AddEvent("Run failed", trace.WithAttributes(
		attribute.String("http_status", http.StatusText(httpStatus)),
		attribute.String("stderr", stderr),
		attribute.String("error", err.Error()),
	))

	json.NewEncoder(w).Encode(resp)
}

func returnSuccessResponse(logger *slog.Logger, tracer trace.Tracer, span trace.Span, w http.ResponseWriter, stdout, stderr string, cid publisher.ResponseCID) {

	defer span.End()

	resp := Response{
		SourceCID: cid.Source,
		OutputCID: cid.Stdout,
		Status:    "completed",
		Stdout:    stdout,
		Stderr:    stderr,
	}

	httpStatus := http.StatusOK
	w.WriteHeader(httpStatus)
	logger.Info("Returning successful response", slog.String("http_status", http.StatusText(httpStatus)), slog.String("source_cid", cid.Source), slog.String("output_cid", cid.Stdout))

	span.AddEvent("Run succeeded", trace.WithAttributes(
		attribute.String("http_status", http.StatusText(httpStatus)),
		attribute.String("stdout", stdout),
		attribute.String("stderr", stderr),
		attribute.String("source_cid", cid.Source),
		attribute.String("output_cid", cid.Stdout),
	))

	json.NewEncoder(w).Encode(resp)
}

func PrometheusMetricsHandler() http.HandlerFunc {

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	}
}

func RunHandler(b *backend.Backend) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var request Request
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			b.Logger.Error("Received invalid HTTP request", slog.Any("HTTPStatus", http.StatusBadRequest), slog.Any("error", err))
			return
		}

		b.Logger.Info("Received HTTP request at /api/run", slog.String("language", request.Language), slog.Int("code_length", len(request.Code)))

		w.Header().Set("Content-Type", "application/json")

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
		defer cancel()

		ctx, span := b.Tracer.Start(ctx, "backend.request", trace.WithAttributes(
			attribute.String("language", request.Language),
		))

		var run backend.RunSpecs = backend.RunSpecs{
			Language: request.Language,
			Source:   []byte(request.Code),
			Start:    time.Now(),
			Span:     span,
		}

		stdout, stderr, responseCID, err := b.Run(ctx, run)
		if err != nil {
			returnFailedResponse(b.Logger, b.Tracer, w, stderr, err)
		} else {
			returnSuccessResponse(b.Logger, b.Tracer, span, w, stdout, stderr, responseCID)
		}
	}
}

func PublishHandler(b *backend.Backend) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// * DEBUG
		log.Printf("Publish request incoming from %s", r.RemoteAddr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	}
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lgtm/internal/api"
	"lgtm/internal/backend"
	"lgtm/internal/compiler"
	"lgtm/internal/executor"
	"lgtm/internal/publisher"
	"lgtm/internal/telemetry"
)

// The local cors function is a  middleware that checks the origin of the request
// and sets the appropriate CORS headers to allow cross-origin requests from the client only.
func cors(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if os.Getenv("ENV") == "production" {
			w.Header().Set("Access-Control-Allow-Origin", "http://lgtm.local")
		} else {
			// In development, allow requests from any origin for testing purposes.
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func newServer(b *backend.Backend) *http.Server {

	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", api.HealthHandler)
	mux.HandleFunc("/api/run", api.RunHandler(b))
	mux.HandleFunc("/api/publish", api.PublishHandler(b))

	return &http.Server{
		Addr:    ":4242",
		Handler: cors(mux),
	}
}

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tracer, shutdownTracing, err := telemetry.InitTracing(ctx, "otel-collector:4317")
	if err != nil {
		log.Fatal(err)
	}
	defer shutdownTracing(context.Background())

	meter, shutdownMetrics, err := telemetry.InitMetrics(ctx, "otel-collector:4317")
	if err != nil {
		log.Fatal(err)
	}
	defer shutdownMetrics(context.Background())

	logger, shutdownLogging, err := telemetry.InitLogging(ctx, "otel-collector:4317")
	if err != nil {
		log.Fatal(err)
	}
	defer shutdownLogging(context.Background())

	sb := compiler.NewWazeroSandbox(tracer, logger)
	exe := executor.NewWazeroExecutor(context.Background(), tracer, logger)
	p := publisher.NewIPFSPublisher(tracer, logger, "http://ipfs:5001")

	b := backend.NewBackend(sb, exe, p, tracer, meter, logger)

	httpserver := newServer(b)

	go func() {
		log.Printf("LGTM Backend server running at http://localhost:4242")
		b.Logger.InfoContext(context.Background(), "LGTM Backend server running at http://localhost:4242")
		log.Fatal(httpserver.ListenAndServe())
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	b.Logger.InfoContext(context.Background(), "Shutting down server...")
	if err := httpserver.Shutdown(shutdownCtx); err != nil {
		b.Logger.ErrorContext(context.Background(), "server shutdown", "error", err)
	}
}

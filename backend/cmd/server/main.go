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

	mux.HandleFunc("/metrics", api.PrometheusMetricsHandler())
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

	shutdownTracing, err := telemetry.InitTracing(ctx, "otel-collector:4317")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdownTracing(context.Background()); err != nil {
			log.Printf("tracing shutdown: %v", err)
		}
	}()

	sb := compiler.NewWazeroSandbox()
	exe := executor.NewWazeroExecutor(context.Background())

	gatewayURL := "http://ipfs:5001"
	p, err := publisher.NewIPFSPublisher(gatewayURL)

	if err != nil {
		log.Fatal(err)
	}

	b := backend.NewBackend(sb, exe, p)

	httpserver := newServer(b)

	go func() {
		log.Println("LGTM Backend server running at http://localhost:4242")
		log.Fatal(httpserver.ListenAndServe())
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := httpserver.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"lgtm/internal/api"
	"lgtm/internal/sandbox"
	"lgtm/internal/server"
)

// The local cors function is a  middleware that checks the origin of the request
// and sets the appropriate CORS headers to allow cross-origin requests from the client only.
func cors(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func newServer(sb *sandbox.Sandbox) *http.Server {

	s := &server.Server{
		Sandbox: sb,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", api.HealthHandler)
	mux.HandleFunc("/api/run", s.RunRequestHandler)
	// mux.HandleFunc("/api/publish", sandbox.PublishRequestHandler)

	return &http.Server{
		Addr:    ":4242",
		Handler: cors(mux),
	}
}

func main() {

	sb := sandbox.NewSandbox(
		64*1024*1024,     // 64 MB memory limit
		10*time.Second,   // 10 seconds timeout
		1024*1024,        // 1 MB max stdout
		1024*1024,        // 1 MB max stderr
		[]string{"/tmp"}, // Allowed directories
	)

	httpserver := newServer(sb)

	fmt.Println("LGTM Backend server running at http://localhost:4242")
	log.Fatal(httpserver.ListenAndServe())

}

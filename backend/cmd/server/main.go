package main

import (
	"fmt"
	"log"
	"net/http"

	"lgtm/internal/api"
)

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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", api.HealthHandler)
	mux.HandleFunc("/api/run", api.RunHandler)

	fmt.Println("LGTM Backend server running at http://localhost:4242")
	log.Fatal(http.ListenAndServe(":4242", cors(mux)))
}

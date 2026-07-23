package api

import (
	"context"
	"encoding/json"
	"fmt"
	"lgtm/internal/backend"
	"net/http"
	"time"
)

type Request struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type Response struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`
}

func RunHandler(b *backend.Backend) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var request Request
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		stdout, stderr, err := b.Run(ctx, []byte(request.Code))
		if err != nil {
			http.Error(w, "Failed to run WASM", http.StatusInternalServerError)
			return
		}

		time.Sleep(2 * time.Second) // Simulate some processing time

		resp := Response{
			ID:     "some-id",
			Status: "completed",
			Stdout: stdout,
			Stderr: stderr,
		}

		json.NewEncoder(w).Encode(resp)
	}

}

func PublishHandler(b *backend.Backend) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// * DEBUG
		fmt.Println("Publish request incoming from ", r.RemoteAddr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	}
}

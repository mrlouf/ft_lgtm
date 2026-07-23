package api

import (
	"context"
	"encoding/json"
	"lgtm/internal/backend"
	"log"
	"net/http"
	"time"
)

type Request struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type Response struct {
	CID    string `json:"cid"`
	Status string `json:"status"`
	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`
}

func returnFailedResponse(w http.ResponseWriter, stderr string, err error) {

	log.Printf("Run failed: %v", err)

	w.WriteHeader(http.StatusInternalServerError)

	resp := Response{
		Status: "failed",
		Stderr: stderr,
	}

	json.NewEncoder(w).Encode(resp)
}

func returnSuccessResponse(w http.ResponseWriter, stdout, stderr, cid string) {

	log.Printf("Run succeeded:\n stdout: %s\n stderr: %s\n cid: %s", stdout, stderr, cid)

	w.WriteHeader(http.StatusOK)
	resp := Response{
		CID:    cid,
		Status: "completed",
		Stdout: stdout,
		Stderr: stderr,
	}

	json.NewEncoder(w).Encode(resp)
}

func RunHandler(b *backend.Backend) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var request Request
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 6*time.Second)
		defer cancel()

		source := []byte(request.Code)
		language := request.Language

		stdout, stderr, cid, err := b.Run(ctx, source, language)
		if err != nil {

			returnFailedResponse(w, stderr, err)

		} else {

			returnSuccessResponse(w, stdout, stderr, cid)

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

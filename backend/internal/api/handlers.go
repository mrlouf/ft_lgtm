package api

import (
	"context"
	"encoding/json"
	"lgtm/internal/backend"
	"lgtm/internal/ipfs"
	"log"
	"net/http"
	"strings"
	"time"
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

func returnFailedResponse(w http.ResponseWriter, stderr string, err error) {

	log.Printf("Run failed: %v", err)

	httpStatus := getHTTPStatusFromError(err)
	w.WriteHeader(httpStatus)

	resp := Response{
		Status: "failed",
		Stderr: stderr,
		Error:  err.Error(),
	}

	json.NewEncoder(w).Encode(resp)
}

func returnSuccessResponse(w http.ResponseWriter, stdout, stderr string, cid ipfs.ResponseCID) {

	log.Printf("Run succeeded:\n stdout: %s\n stderr: %s\n source cid: %s\n output cid: %s", stdout, stderr, cid.Source, cid.Stdout)

	w.WriteHeader(http.StatusOK)
	resp := Response{
		SourceCID: cid.Source,
		OutputCID: cid.Stdout,
		Status:    "completed",
		Stdout:    stdout,
		Stderr:    stderr,
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
		ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
		defer cancel()

		source := []byte(request.Code)
		language := request.Language

		stdout, stderr, responseCID, err := b.Run(ctx, source, language)
		if err != nil {
			returnFailedResponse(w, stderr, err)
		} else {
			returnSuccessResponse(w, stdout, stderr, responseCID)
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

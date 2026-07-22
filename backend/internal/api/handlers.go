package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Health check request from ", r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func runRequest(ctx context.Context) (RunResponse, error) {

	var response RunResponse

	select {
	case <-ctx.Done():
		fmt.Println("Request canceled by the client")
		return response, fmt.Errorf("request canceled")
	case <-time.After(5 * time.Second):
		// Simulate job processing
		response.Status = "completed"
		fmt.Println("Job completed:", response)
		return response, nil
	default:
		// Simulate job processing
		time.Sleep(5 * time.Second)
		response.Status = "completed"
		fmt.Println("Job completed:", response)
	}

	return response, nil

}

func RunRequestHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Run Request incoming from ", r.RemoteAddr)

	var request RunRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// * DEBUG
	fmt.Println("Received RunRequest:", request)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := runRequest(ctx)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

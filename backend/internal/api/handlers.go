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

func runRequest(ctx context.Context, request RunRequest, resp *RunResponse) {

	// * DEBUG
	fmt.Println("Running request with language:", request.Language)
	fmt.Println("Running request with snippet:", request.Code)

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Task completed successfully")
		resp.Status = "completed"
	case <-ctx.Done():
		fmt.Println("Task canceled due to timeout")
		resp.Status = "timeout"
	}
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

	var resp RunResponse
	go runRequest(ctx, request, &resp)

	<-ctx.Done()

	json.NewEncoder(w).Encode(resp)
}

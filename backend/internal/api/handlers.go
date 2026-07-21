package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Health check request from ", r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func RunHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Run Request incoming from ", r.RemoteAddr)

	var request RunRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Received RunRequest:", request)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := RunResponse{
		ID:     "job-123",
		Status: "queued",
	}

	json.NewEncoder(w).Encode(response)
}

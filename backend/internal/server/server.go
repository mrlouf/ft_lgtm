package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"lgtm/internal/api"
	"lgtm/internal/sandbox"
)

type Server struct {
	Sandbox *sandbox.Sandbox
}

func (s *Server) RunRequestHandler(w http.ResponseWriter, r *http.Request) {

	// * DEBUG
	fmt.Println("Request incoming from ", r.RemoteAddr)

	var request api.Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// * DEBUG
	fmt.Println("Received Request:", request)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp := api.RunRequest(ctx, request)

	json.NewEncoder(w).Encode(resp)

}

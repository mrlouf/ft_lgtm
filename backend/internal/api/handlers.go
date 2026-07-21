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

	// var request RunRequest

	var request RunRequest

	json.NewDecoder(r.Body).Decode(&request)

	fmt.Println("Request: ", request)

}

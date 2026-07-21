package api

import (
	"fmt"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Health check request from ", r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func RunHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("You've reached the run handler")

}

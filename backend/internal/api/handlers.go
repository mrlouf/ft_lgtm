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

	fmt.Println("Run Request incoming from ", r.RemoteAddr)

	// var request RunRequest

	body := r.Body

	fmt.Println(body)

}

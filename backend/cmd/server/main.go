package main

import (
	"fmt"
	"log"
	"net/http"

	"lgtm/internal/api"
)

func main() {

	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/run", api.RunHandler)

	fmt.Println("LGTM Backend server running at http://localhost:4242")

	log.Fatal(http.ListenAndServe(":4242", nil))

}

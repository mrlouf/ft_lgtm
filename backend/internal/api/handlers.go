package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Health check request from ", r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func RunRequest(ctx context.Context, request Request) Response {

	// * DEBUG
	fmt.Println("Running request with language:", request.Language)
	fmt.Println("Running request with snippet:", request.Code)

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Task completed successfully")
		return Response{Status: "completed"}
	case <-ctx.Done():
		fmt.Println("Task canceled")
		return Response{Status: "canceled"}
	}
}

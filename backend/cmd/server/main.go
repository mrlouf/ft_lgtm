package main

import (
	"fmt"
	"log"
	"net/http"

	"lgtm/internal/api"
	"lgtm/internal/backend"
	"lgtm/internal/ipfs"
	"lgtm/internal/sandbox"
)

// The local cors function is a  middleware that checks the origin of the request
// and sets the appropriate CORS headers to allow cross-origin requests from the client only.
func cors(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func newServer(sb *sandbox.WazeroSandbox, exe *sandbox.WazeroExecutor, ipfs *ipfs.IPFS) *http.Server {

	b := &backend.Backend{
		Compiler:  sb,
		Executor:  exe,
		Publisher: ipfs,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", api.HealthHandler)
	mux.HandleFunc("/api/run", api.RunHandler(b))
	mux.HandleFunc("/api/publish", api.PublishHandler(b))

	return &http.Server{
		Addr:    ":4242",
		Handler: cors(mux),
	}
}

func main() {

	sb := sandbox.NewWazeroSandbox()
	exe := sandbox.NewWazeroExecutor(sb)
	ipfs := ipfs.NewIPFSClient()

	httpserver := newServer(sb, exe, ipfs)

	fmt.Println("LGTM Backend server running at http://localhost:4242")
	log.Fatal(httpserver.ListenAndServe())

}

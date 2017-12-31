package frontend

import (
	"log"
	"net/http"
)

// Serve the web browser frontend
func Serve() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/subscribe", subscribe)
	mux.HandleFunc("/pub", getPublicKey)
	mux.Handle("/", http.FileServer(http.Dir("./frontend")))

	server := &http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
	}

	return server.ListenAndServe()
}

func register(w http.ResponseWriter, r *http.Request) {
	log.Println(".")
}

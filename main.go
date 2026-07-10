package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	fsHandler := http.FileServer(http.Dir("."))
	cfg := &apiConfig{}

	mux.Handle("/app/", cfg.middlewareMetricsInc(fsHandler))
	mux.HandleFunc("/healthz", cfg.handlerHealth)
	mux.HandleFunc("/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/reset", cfg.handlerReset)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server starting on http://localhost:8080")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

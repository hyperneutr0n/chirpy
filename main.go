package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const filePathRoot = "app"
	const port = "8080"

	mux := http.NewServeMux()
	apiCfg := apiConfig{fsHit: atomic.Int32{}}

	fs := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port: %v\n", port)
	log.Fatal(server.ListenAndServe())
}


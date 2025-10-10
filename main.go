package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fsHit atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fsHit.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	count := cfg.fsHit.Load()

	s := fmt.Sprintf(
		`<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`, count)
	data := []byte(s)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	if err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fsHit.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

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

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body string `json:"body"`
	}
	
	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("error reading body: %v", err)
		return
	}

	length := len(req.Body)
	if length > 140 {
		sendError(w, 400, "Chirp is too long")
		return;
	}

	words := strings.Split(req.Body, " ")
	for i, word := range words{
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "formax" {
			words[i] = "****"
		}
	}

	respBody := struct {
		CleanedBody string `json:"cleaned_body"`
	}{CleanedBody: strings.Join(words, " ")}
	sendJSON(w, 200, respBody)
}

func sendError(w http.ResponseWriter, code int, message string) {
	type respJson struct {
		Error string `json:"error"`
	}
	respBody := respJson{
		Error: message,
	}
	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("error when marshalling response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("error when marshalling response: %v", err)
		return
	}
}

func sendJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error when marshalling response: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("error when marshalling response: %v", err)
		return
	}
}
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hyperneutr0n/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("error reading body: %v", err)
		sendError(w, 500, "Failed creating user")
		return
	}

	length := len(req.Body)
	if length > 140 {
		sendError(w, 400, "Chirp is too long")
		return
	}

	words := strings.Split(req.Body, " ")
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "formax" {
			words[i] = "****"
		}
	}

	newChirp := database.CreateChirpParams{
		Body: strings.Join(words, " "),
		UserID: req.UserID,
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), newChirp)
	if err != nil {
		log.Printf("failed creating chirp: %v", err)
		return
	}

	sendJSON(w, 201, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("error retrieving chirps: %v", err)
		return
	}

	sendJSON(w, 200, chirps)
}
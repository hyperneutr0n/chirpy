package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hyperneutr0n/chirpy/internal/auth"
	"github.com/hyperneutr0n/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error when getting bearer token from header: %v", err)
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println(bearer)
	loggedInUserID, err := auth.ValidateJWT(bearer, cfg.secret)
	if err != nil {
		log.Printf("error when validating bearer token: %v", err)
		sendError(w, http.StatusUnauthorized, "Invalid bearer token")
		return
	}

	type request struct {
		Body   string    `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err = decoder.Decode(&req)
	if err != nil {
		log.Printf("error reading body: %v", err)
		sendError(w, 500, "Failed creating chirp")
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
		Body:   strings.Join(words, " "),
		UserID: loggedInUserID,
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), newChirp)
	if err != nil {
		log.Printf("failed creating chirp: %v", err)
		sendError(w, 500, "Failed creating chirp")
		return
	}

	sendJSON(w, 201, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("error retrieving chirps: %v", err)
		sendError(w, 500, "Failed retrieving chirps")
		return
	}

	sendJSON(w, 200, chirps)
}

func (cfg *apiConfig) handlerFindChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	_chirpID, err := uuid.Parse(chirpID)
	if err != nil {
		log.Printf("failed parsing uuid: %v", err)
		sendError(w, 400, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.db.FindChirp(r.Context(), _chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, 404, "Chirp not found")
			return
		}
		log.Printf("failed finding chirp: %v", err)
		sendError(w, 500, "Internal server error")
		return
	}

	sendJSON(w, 200, chirp)
}

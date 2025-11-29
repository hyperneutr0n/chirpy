package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hyperneutr0n/chirpy/internal/auth"
	"github.com/hyperneutr0n/chirpy/internal/database"
)

const DefaultExpiresIn = 1 * time.Hour

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerRegister(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("error reading body: %v", err)
		sendError(w, 400, "Invalid request body")
		return
	}

	password, err := auth.HashPasword(req.Password)
	if err != nil {
		log.Printf("error when trying to hash password %v", err)
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    req.Email,
		Password: password,
	})
	if err != nil {
		log.Printf("error creating user: %v", err)
		sendError(w, 500, "Failed creating user")
		return
	}

	sendJSON(w, 201, user)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("error reading body: %v", err)
		sendError(w, 400, "Invalid request body")
		return
	}

	user, err := cfg.db.LoginUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("error fetching user data: %v", err)
		sendError(w, 500, "Internal server error")
		return
	}

	match, err := auth.CheckHash(req.Password, user.Password)
	if err != nil {
		log.Printf("error when checking hash: %v", err)
		sendError(w, 500, "Internal server error")
		return
	}

	if !match {
		sendError(w, 401, "Incorrect email or password")
		return
	}

	loggedInUser, err := cfg.db.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		log.Printf("error when fetching logged in user: %v", err)
		sendError(w, 500, "Internal server error")
		return
	}

	expiresIn := DefaultExpiresIn
	if req.ExpiresInSeconds < 0 || req.ExpiresInSeconds > 3600 {
		expiresIn = time.Duration(req.ExpiresInSeconds)
	}

	token, err := auth.MakeJWT(loggedInUser.ID, cfg.secret, expiresIn)
	if err != nil {
		log.Printf("error when creating jwt: %v", err)
		sendError(w, 500, "Internal server error")
		return
	}

	sendJSON(w, 200, User{
		ID:        loggedInUser.ID,
		Email:     loggedInUser.Email,
		CreatedAt: loggedInUser.CreatedAt,
		UpdatedAt: loggedInUser.UpdatedAt,
		Token:     token,
	})
}

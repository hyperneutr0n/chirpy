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
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	token, err := auth.MakeJWT(loggedInUser.ID, cfg.secret, DefaultExpiresIn)
	if err != nil {
		log.Printf("error when creating jwt: %v", err)
		sendError(w, 500, "Internal server error")
		return
	}

	refreshToken, _ := auth.MakeRefreshToken()
	_refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    loggedInUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour).UTC(),
	})
	if err != nil {
		log.Printf("error when creating refresh token: %v", err)
		sendError(w, 500, "Internal server error")
		return
	}

	sendJSON(w, 200, User{
		ID:           loggedInUser.ID,
		Email:        loggedInUser.Email,
		IsChirpyRed:  loggedInUser.IsChirpyRed,
		CreatedAt:    loggedInUser.CreatedAt,
		UpdatedAt:    loggedInUser.UpdatedAt,
		Token:        token,
		RefreshToken: _refreshToken.Token,
	})
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error getting refresh token from header: %v", err)
		sendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tokenRecord, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("error getting refresh token: %v", err)
		sendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if time.Now().UTC().After(tokenRecord.ExpiresAt) {
		sendError(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	if tokenRecord.RevokedAt.Valid {
		sendError(w, http.StatusUnauthorized, "Refresh token revoked")
		return
	}

	token, err := auth.MakeJWT(tokenRecord.UserID, cfg.secret, DefaultExpiresIn)
	if err != nil {
		log.Printf("error creating jwt: %v", err)
		sendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	sendJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error getting refresh token from header: %v", err)
		sendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("error when revoking refresh token: %v", err)
		sendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	sendJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error when getting bearer token from header: %v", err)
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	loggedInUserID, err := auth.ValidateJWT(bearer, cfg.secret)
	if err != nil {
		log.Printf("error when validating bearer token: %v", err)
		sendError(w, http.StatusUnauthorized, "Invalid bearer token")
		return
	}

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err = decoder.Decode(&req)
	if err != nil {
		log.Printf("error reading body: %v", err)
		sendError(w, 400, "Invalid request body")
		return
	}

	password, err := auth.HashPasword(req.Password)
	if err != nil {
		log.Printf("error when trying to hash password %v", err)
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:    req.Email,
		Password: password,
		ID:       loggedInUserID,
	})
	if err != nil {
		log.Printf("error creating user: %v", err)
		sendError(w, 500, "Failed creating user")
		return
	}

	sendJSON(w, 200, user)
}

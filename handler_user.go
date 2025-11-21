package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID			uuid.UUID	`json:"id"`
	Email		string		`json:"email"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

func (cfg *apiConfig) handlerRegister(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("error reading body: %v", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("error creating user: %v", err)
		sendError(w, 500, "Failed creating user")
		return
	}

	responseUser := User{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	sendJSON(w, 201, responseUser)
}
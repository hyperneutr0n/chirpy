package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hyperneutr0n/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		sendError(w, 401, err.Error())
		return
	}
	if apiKey != cfg.polkaKey {
		sendError(w, 401, "Polka key doesn't match");
		return
	}

	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err = decoder.Decode(&req)
	if err != nil {
		log.Printf("error reading body: %v", err)
		sendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if req.Event != "user.upgraded" {
		sendJSON(w, 204, "")
		return
	}

	err = cfg.db.UpgradeUser(r.Context(), req.Data.UserID)
	if err != nil {
		log.Printf("failed upgrading user: %v", err)
		sendError(w, http.StatusNotFound, "User not found")
		return
	}
	
	sendJSON(w, 204, "")
}

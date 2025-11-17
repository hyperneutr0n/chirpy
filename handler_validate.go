package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

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
		return
	}

	words := strings.Split(req.Body, " ")
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "formax" {
			words[i] = "****"
		}
	}

	respBody := struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: strings.Join(words, " "),
	}
	sendJSON(w, 200, respBody)
}

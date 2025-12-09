package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func sendJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error when marshalling response: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if code != http.StatusNoContent {
		_, err = w.Write(data)
		if err != nil {
			log.Printf("error when writing response: %v", err)
			return
		}
	}
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
		log.Printf("error when writing response: %v", err)
		return
	}
}

package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	
	mux.Handle("/", http.FileServer(http.Dir(".")))

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Println("Starting server at localhost:8080")
	log.Fatalln(server.ListenAndServe())
}
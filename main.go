package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/hyperneutr0n/chirpy/internal/database"
)

func main() {
	const filePathRoot = "app"
	const port = "8080"

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	apiCfg := &apiConfig{
		fsHit: atomic.Int32{},
		db: dbQueries,
	}

	fs := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("GET /api/healthz", handlerHealthz)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/users", apiCfg.handlerRegister)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port: %v\n", port)
	log.Fatal(server.ListenAndServe())
}


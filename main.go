package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
	"embed"

	"github.com/pressly/goose/v3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/hyperneutr0n/chirpy/internal/database"
)

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

func main() {
	const filePathRoot = "app"
	const port = "8080"

	if err := godotenv.Load(); err != nil {
		log.Printf("No env file found")
	}

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
        log.Fatal(err)
    }

    goose.SetBaseFS(embedMigrations)

    log.Println("Running database migrations...")
    if err := goose.Up(db, "sql/schema"); err != nil {
        log.Fatalf("Error running migrations: %v", err)
    }
    log.Println("Migrations completed successfully!")

	dbQueries := database.New(db)

	secret := os.Getenv("SECRET")

	polkaKey := os.Getenv("POLKA_KEY")

	apiCfg := &apiConfig{
		fsHit:    atomic.Int32{},
		db:       dbQueries,
		secret:   secret,
		polkaKey: polkaKey,
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("GET /api/healthz", handlerHealthz)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/users", apiCfg.handlerRegister)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerFindChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhooks)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting Chirpy on port: %v\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	<-stop
	log.Println("Shutting down Chirpy...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown failed: %v", err)
	}

	log.Println("Chirpy exited cleanly.")
}

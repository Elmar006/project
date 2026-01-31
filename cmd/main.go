package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Elmar006/project/internal/api"
	"github.com/Elmar006/project/internal/db"
	"github.com/Elmar006/project/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	web := "../web"
	logger.Init()

	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: Could not load .env file %v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "7540"
		log.Printf("PORT environment variable not set, using default: %s", port)
	}
	dbAddr := os.Getenv("TODO_DBFILE")
	if dbAddr == "" {
		dbAddr = "scheduler.db"
	}
	err := db.Init(dbAddr)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir(web)))
	log.Printf("Starting server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

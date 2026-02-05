package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Elmar006/project/internal/api"
	"github.com/Elmar006/project/internal/db"
	"github.com/Elmar006/project/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	logger.Init()

	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	projectDir := curDir
	if filepath.Base(curDir) == "cmd" {
		projectDir = filepath.Dir(curDir)
	}

	envPath := filepath.Join(projectDir, ".env")
	if err := godotenv.Load(envPath); err != nil {
		logger.L().Errorf("Warning: Could not load .env file %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7540"
		logger.L().Infof("The default port is used: %s", port)
	}

	dbAddr := os.Getenv("TODO_DBFILE")
	if dbAddr == "" {
		dbAddr = "scheduler.db"
		logger.L().Infof("The default database address is used: %s", dbAddr)
	}

	password := os.Getenv("TODO_PASSWORD")
	if !filepath.IsAbs(dbAddr) {
		dbAddr = filepath.Join(projectDir, dbAddr)
	}

	err = db.Init(dbAddr)
	if err != nil {
		logger.L().Errorf("Failed to initialize database: %v", err)
		os.Exit(1)
	}
	defer db.DB.Close()

	api.PasswordCheck(password)
	api.Init()
	webDir := filepath.Join(projectDir, "web")

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	logger.L().Infof("Starting server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.L().Error("Failed to start server")
		os.Exit(1)
	}
}

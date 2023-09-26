package config

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultPort      = "8080"
	defaultJwtSecret = "aSecret"
)

var (
	PORT         string
	DATABASE_DSN string
	JWT_SECRET   string
)

func LoadVariables() {
	log.SetFlags(log.Llongfile + log.LstdFlags)
	if _, err := os.Stat(".env"); err == nil {
		slog.Info("Load .env file")
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file: ", err)
		}
	}

	// Environment Variables
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = defaultPort
	}

	DATABASE_DSN = os.Getenv("DATABASE_DSN")
	if DATABASE_DSN == "" {
		log.Fatal("Environment variable DATABASE_DSN is mandatory")
	}

	JWT_SECRET = os.Getenv("JWT_SECRET")
	if JWT_SECRET == "" {
		JWT_SECRET = defaultJwtSecret
	}
}

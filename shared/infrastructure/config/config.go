package config

import (
	golog "log"
	"os"

	"github.com/gofiber/fiber/v2/log"

	"github.com/joho/godotenv"
)

const (
	defaultPort      = "8080"
	defaultJwtSecret = "aSecret"
)

var (
	PORT         string
	DATABASE_DSN string
	JWT_SECRET   = []byte(defaultJwtSecret)
)

func LoadVariables() {
	golog.SetFlags(golog.Llongfile + golog.LstdFlags)
	if _, err := os.Stat(".env"); err == nil {
		log.Info("Load .env file")
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

	if env := os.Getenv("JWT_SECRET"); env != "" {
		JWT_SECRET = []byte(env)
	}
}

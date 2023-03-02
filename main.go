package main

import (
	"api-your-accounts/shared/infrastructure"
	"api-your-accounts/shared/infrastructure/db"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	osStat         = os.Stat
	dotenvLoad     = godotenv.Load
	logFatal       = log.Fatal
	newDB          = db.NewDB
	newMongoClient = db.NewMongoClient
	newServer      = infrastructure.NewServer
)

func main() {
	log.SetFlags(log.Lshortfile + log.LUTC + log.LstdFlags)
	if _, err := osStat(".env"); err == nil {
		log.Println("Load .env file")
		err = dotenvLoad()
		if err != nil {
			logFatal("Error loading .env file: ", err)
		}
	}

	newDB()
	newMongoClient()
	newServer()
}

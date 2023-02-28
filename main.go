package main

import (
	"api-your-accounts/shared/infrastructure"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	osStat     = os.Stat
	dotenvLoad = godotenv.Load
	logFatal   = log.Fatal
	newDB      = infrastructure.NewDB
	newServer  = infrastructure.NewServer
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

	db := newDB()
	newServer(db)
}

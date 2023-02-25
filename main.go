package main

import (
	"api-your-accounts/shared/infrastructure"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.Lshortfile + log.LUTC + log.LstdFlags)
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Load .env file")
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file: ", err)
		}
	}

	db := infrastructure.NewDB()
	infrastructure.NewServer(db)
}

package main

import (
	"api-your-accounts/infrastructure"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Load .env file")
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file: ", err)
		}
	}

	infrastructure.NewServer()
}

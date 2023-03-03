package main

import (
	"api-your-accounts/shared/infrastructure"
	"api-your-accounts/shared/infrastructure/db"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "api-your-accounts/docs"
)

var (
	osStat         = os.Stat
	dotenvLoad     = godotenv.Load
	logFatal       = log.Fatal
	newDB          = db.NewDB
	newMongoClient = db.NewMongoClient
	newServer      = infrastructure.NewServer
)

// Main godoc
//
//	@title			Your Accounts API
//	@version		1.0
//	@description	This is the API from project Your Accounts
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	Your Accounts Support
//	@contact.email	jonathancastaneda@jsc-developer.me
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@BasePath		/
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

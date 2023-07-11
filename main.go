//go:generate swag fmt
//go:generate swag init

package main

import (
	"log"
	"os"
	"your-accounts-api/shared/infrastructure"
	"your-accounts-api/shared/infrastructure/db"
	"your-accounts-api/shared/infrastructure/handler"
	user "your-accounts-api/user/infrastructure/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "your-accounts-api/docs"
)

var (
	routers = []infrastructure.Router{
		user.NewRoute,
		handler.NewRoute,
	}
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
	log.SetFlags(log.Llongfile + log.LstdFlags)
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Load .env file")
		err = godotenv.Load()
		if err != nil {
			log.Panic("Error loading .env file: ", err)
		}
	}

	db.NewDB()

	// Init server
	server := infrastructure.NewServer(false)
	server.
		AddRoute(infrastructure.Route{
			Method: fiber.MethodGet,
			Path:   "/swagger/*",
			Handler: swagger.New(swagger.Config{
				Title: "Doc API",
			}),
		})

	for _, router := range routers {
		server.AddRoute(router)
	}

	server.Listen()
}

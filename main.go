//go:generate swag fmt
//go:generate swag init

package main

import (
	"api-your-accounts/shared/infrastructure"
	"api-your-accounts/shared/infrastructure/db"
	"api-your-accounts/shared/infrastructure/handler"
	user "api-your-accounts/user/infrastructure/handler"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "api-your-accounts/docs"
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
	var err error
	time.Local, err = time.LoadLocation("America/Bogota")
	if err != nil {
		log.Panic("Error setting location for time: ", err)
	}

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

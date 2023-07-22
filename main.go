//go:generate swag fmt
//go:generate swag init

package main

import (
	"your-accounts-api/shared/infrastructure"
	"your-accounts-api/shared/infrastructure/config"
	"your-accounts-api/shared/infrastructure/db"
	"your-accounts-api/shared/infrastructure/handler"
	"your-accounts-api/shared/infrastructure/injection"
	user "your-accounts-api/user/infrastructure/handler"
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
	config.LoadVariables()
	db.NewDB()
	injection.LoadInstances()

	// Init server
	server := infrastructure.NewServer(false)

	for _, router := range routers {
		server.AddRoute(router)
	}

	server.Listen()
}

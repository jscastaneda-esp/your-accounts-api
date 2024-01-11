//go:generate swag fmt
//go:generate swag init

package main

import (
	"your-accounts-api/shared/infrastructure"
	"your-accounts-api/shared/infrastructure/config"
	"your-accounts-api/shared/infrastructure/db"
	"your-accounts-api/shared/infrastructure/handler"
	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/schedule"
	users "your-accounts-api/users/infrastructure/handler"
)

var (
	routers = []infrastructure.Router{
		users.NewRoute,
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
	server := infrastructure.NewServer(false)

	for _, router := range routers {
		server.AddRoute(router)
	}

	schedule.Start()
	defer schedule.Stop()

	server.Listen()
}

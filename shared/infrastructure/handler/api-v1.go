package handler

import (
	budgets "your-accounts-api/budgets/infrastructure/handler"
	"your-accounts-api/shared/domain/jwt"
	"your-accounts-api/shared/infrastructure/config"
	logs "your-accounts-api/shared/infrastructure/handler/logs"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

var (
	logsRouter    = logs.NewRoute
	budgetsRouter = budgets.NewRoute
)

func NewRoute(app fiber.Router) {
	api := app.Group("/api/v1")
	// Middleware
	{
		api.Use(jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(config.JWT_SECRET)},
			Claims:     new(jwt.JwtUserClaims),
		}))
	}

	// Routes
	{
		logsRouter(api)
		budgetsRouter(api)
	}
}

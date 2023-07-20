package handler

import (
	budget "your-accounts-api/budget/infrastructure/handler"
	project "your-accounts-api/project/infrastructure/handler"
	"your-accounts-api/shared/domain/jwt"
	"your-accounts-api/shared/infrastructure/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

var (
	projectRouter = project.NewRoute
	budgetRouter  = budget.NewRoute
)

func NewRoute(app fiber.Router) {
	api := app.Group("/api/v1")
	// Middleware
	{
		api.Use(jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(config.JWT_SECRET)},
			Claims:     &jwt.JwtUserClaims{},
		}))
	}

	// Routes
	{
		projectRouter(api)
		budgetRouter(api)
	}
}

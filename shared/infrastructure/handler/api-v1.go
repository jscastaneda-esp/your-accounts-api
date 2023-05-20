package handler

import (
	budget "api-your-accounts/budget/infrastructure/handler"
	project "api-your-accounts/project/infrastructure/handler"
	"os"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/gofiber/jwt/v3"
)

var (
	projectRouter = project.NewRoute
	budgetRouter  = budget.NewRoute
)

func NewRoute(app fiber.Router) {
	jwtSecret := os.Getenv("JWT_SECRET")

	api := app.Group("/api/v1")
	// Middleware
	{
		api.Use(jwt.New(jwt.Config{
			SigningKey: []byte(jwtSecret),
		}))
	}

	// Routes
	{
		projectRouter(api)
		budgetRouter(api)
	}
}

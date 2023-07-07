package handler

import (
	"os"
	budget "your-accounts-api/budget/infrastructure/handler"
	project "your-accounts-api/project/infrastructure/handler"

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

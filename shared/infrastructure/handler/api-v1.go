package handler

import (
	budgets "your-accounts-api/budgets/infrastructure/handler"
	"your-accounts-api/shared/domain"
	"your-accounts-api/shared/infrastructure/config"
	logs "your-accounts-api/shared/infrastructure/handler/logs"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var (
	logsRouter    = logs.NewRoute
	budgetsRouter = budgets.NewRoute
)

func NewRoute(e *echo.Echo) {
	api := e.Group("/api/v1")

	// Middleware
	{
		api.Use(echojwt.WithConfig(echojwt.Config{
			SigningKey: []byte(config.JWT_SECRET),
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(domain.JwtUserClaims)
			},
		}))
	}

	// Routes
	{
		logsRouter(api)
		budgetsRouter(api)
	}
}

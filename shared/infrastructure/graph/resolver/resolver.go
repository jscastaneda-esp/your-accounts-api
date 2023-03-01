//go:generate go run github.com/99designs/gqlgen generate

package resolver

import (
	"api-your-accounts/shared/infrastructure/mongodb"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB          *gorm.DB
	MongoClient *mongodb.Client
}

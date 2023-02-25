//go:generate go run github.com/99designs/gqlgen generate

package resolver

import "gorm.io/gorm"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *gorm.DB
}

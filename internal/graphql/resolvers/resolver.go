package resolvers

import (
	"errors"
	"write-stream-go/internal/auth"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB          *gorm.DB
	AuthService *auth.AuthService
}

var ErrUnauthenticated = errors.New("user not authenticated")

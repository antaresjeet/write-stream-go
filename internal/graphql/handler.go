package graphql

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"write-stream-go/internal/auth"
	"write-stream-go/internal/graphql/generated"
	"write-stream-go/internal/graphql/resolvers"
	"write-stream-go/internal/models"
)

func NewHandler(db *gorm.DB, authService *auth.AuthService) gin.HandlerFunc {
	resolver := &resolvers.Resolver{
		DB:          db,
		AuthService: authService,
	}

	c := generated.Config{Resolvers: resolver}
	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var user models.User
		if err := db.First(&user, "id = ?", userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "user", &user)
		c.Request = c.Request.WithContext(ctx)

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// NewPlaygroundHandler returns a Gin handler function for the GraphQL Playground
func NewPlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/graphql")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

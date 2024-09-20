package main

import (
	"log"

	"write-stream-go/internal/auth"
	"write-stream-go/internal/config"
	"write-stream-go/internal/database"
	"write-stream-go/internal/graphql"
	"write-stream-go/internal/handlers"
	"write-stream-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)

	// Initialize auth service
	authService := auth.NewAuthService(cfg, db)

	// Initialize GraphQL handler
	graphqlHandler := graphql.NewHandler(db, authService)

	// Initialize user handler
	userHandler := handlers.NewUserHandler(db)

	r := gin.Default()

	// Set up middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// Set up routes
	r.POST("/auth/google", authService.GoogleAuthHandler)
	r.GET("/auth/google/callback", authService.GoogleCallbackHandler)

	// User creation endpoint
	r.POST("/users", userHandler.CreateUser)

	// GraphQL endpoint with authentication
	r.POST("/graphql", middleware.JWTAuthMiddleware(cfg.JWTSecret), graphqlHandler)

	// GraphQL Playground
	r.GET("/playground", graphql.NewPlaygroundHandler())

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

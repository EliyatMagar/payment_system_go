package main

import (
	"log"
	"payment-system/config"
	"payment-system/database"
	"payment-system/routes"
	"payment-system/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Validate required configuration
	if cfg.StripeKey == "" {
		log.Fatal("STRIPE_SECRET_KEY is required")
	}

	// Initialize database
	err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Stripe service
	stripeService := services.NewStripeService(cfg)

	// Create Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, stripeService)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

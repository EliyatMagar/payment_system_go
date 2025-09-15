package routes

import (
	"net/http"
	"payment-system/handlers"
	"payment-system/middleware"
	"payment-system/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, stripeService *services.StripeService) {
	// Apply middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.AuthMiddleware())

	paymentHandler := handlers.NewPaymentHandler(stripeService)
	customerHandler := handlers.NewCustomerHandler(stripeService)

	api := router.Group("/api/v1")
	{
		// Customer routes
		api.POST("/customers", customerHandler.CreateCustomer)
		api.GET("/customers", customerHandler.ListCustomers)
		api.GET("/customers/:id", customerHandler.GetCustomer)
		api.DELETE("/customers/:id", customerHandler.DeleteCustomer)

		// Payment routes
		api.POST("/payment-intents", paymentHandler.CreatePaymentIntent)
		api.GET("/payment-intents/:id/status", paymentHandler.GetPaymentStatus)
		api.GET("/payments", paymentHandler.ListPayments)

		// Webhook route (no auth required)
		api.POST("/webhooks", paymentHandler.HandleWebhook)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

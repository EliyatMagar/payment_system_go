package handlers

import (
	"net/http"
	"payment-system/database"
	"payment-system/models"
	"payment-system/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerHandler struct {
	stripeService *services.StripeService
	db            *gorm.DB
}

func NewCustomerHandler(stripeService *services.StripeService) *CustomerHandler {
	return &CustomerHandler{
		stripeService: stripeService,
		db:            database.GetDB(),
	}
}

type CreateCustomerRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if customer already exists
	var existingCustomer models.Customer
	if err := h.db.Where("email = ?", req.Email).First(&existingCustomer).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Customer already exists"})
		return
	}

	// Create customer in Stripe
	stripeCustomer, err := h.stripeService.CreateCustomer(req.Email, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save customer to database
	customer := models.Customer{
		StripeID: stripeCustomer.ID,
		Email:    req.Email,
		Name:     req.Name,
	}

	if err := h.db.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        customer.ID,
		"stripe_id": customer.StripeID,
		"email":     customer.Email,
		"name":      customer.Name,
	})
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")

	var customer models.Customer
	if err := h.db.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) ListCustomers(c *gin.Context) {
	var customers []models.Customer
	if err := h.db.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	var customer models.Customer
	if err := h.db.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	if err := h.db.Delete(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

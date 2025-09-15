package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"payment-system/database"
	"payment-system/models"
	"payment-system/services"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	stripeService *services.StripeService
	db            *gorm.DB
}

func NewPaymentHandler(stripeService *services.StripeService) *PaymentHandler {
	return &PaymentHandler{
		stripeService: stripeService,
		db:            database.GetDB(),
	}
}

type CreatePaymentIntentRequest struct {
	Amount      int64             `json:"amount" binding:"required,min=50"`
	Currency    string            `json:"currency" binding:"required"`
	CustomerID  uint              `json:"customer_id" binding:"required"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var req CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get customer from database
	var customer models.Customer
	if err := h.db.First(&customer, req.CustomerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Convert map[string]string to map[string]interface{} for Stripe
	stripeMetadata := make(map[string]string)
	for k, v := range req.Metadata {
		stripeMetadata[k] = v
	}

	// Create payment intent in Stripe
	pi, err := h.stripeService.CreatePaymentIntent(
		req.Amount,
		req.Currency,
		customer.StripeID,
		stripeMetadata,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert map[string]string to models.JSONB for database
	jsonbMetadata := make(models.JSONB)
	for k, v := range req.Metadata {
		jsonbMetadata[k] = v
	}

	// Save payment intent to database
	paymentIntent := models.PaymentIntent{
		StripeID:    pi.ID,
		CustomerID:  customer.ID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Status:      string(pi.Status),
		Description: req.Description,
		Metadata:    jsonbMetadata,
	}

	if err := h.db.Create(&paymentIntent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"payment_intent_id": pi.ID,
		"client_secret":     pi.ClientSecret,
		"status":            pi.Status,
		"amount":            pi.Amount,
		"currency":          pi.Currency,
	})
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	signature := c.GetHeader("Stripe-Signature")
	event, err := h.stripeService.HandleWebhook(payload, signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.handlePaymentIntentSucceeded(&paymentIntent)

	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.handlePaymentIntentFailed(&paymentIntent)

	case "payment_intent.canceled":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.handlePaymentIntentCanceled(&paymentIntent)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *PaymentHandler) handlePaymentIntentSucceeded(pi *stripe.PaymentIntent) {
	// Update payment intent status in database
	h.db.Model(&models.PaymentIntent{}).
		Where("stripe_id = ?", pi.ID).
		Update("status", string(pi.Status))

	// Get payment intent from database to get the ID
	var paymentIntent models.PaymentIntent
	h.db.Where("stripe_id = ?", pi.ID).First(&paymentIntent)

	// Create payment record
	payment := models.Payment{
		StripeID:        pi.ID,
		PaymentIntentID: paymentIntent.ID,
		Amount:          pi.Amount,
		Currency:        string(pi.Currency),
		Status:          string(pi.Status),
	}
	h.db.Create(&payment)
}

func (h *PaymentHandler) handlePaymentIntentFailed(pi *stripe.PaymentIntent) {
	h.db.Model(&models.PaymentIntent{}).
		Where("stripe_id = ?", pi.ID).
		Update("status", string(pi.Status))
}

func (h *PaymentHandler) handlePaymentIntentCanceled(pi *stripe.PaymentIntent) {
	h.db.Model(&models.PaymentIntent{}).
		Where("stripe_id = ?", pi.ID).
		Update("status", string(pi.Status))
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	paymentIntentID := c.Param("id")

	pi, err := h.stripeService.GetPaymentIntent(paymentIntentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment intent not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   pi.Status,
		"amount":   pi.Amount,
		"currency": pi.Currency,
	})
}

func (h *PaymentHandler) ListPayments(c *gin.Context) {
	var payments []models.Payment
	if err := h.db.Preload("PaymentIntent").Preload("PaymentIntent.Customer").Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

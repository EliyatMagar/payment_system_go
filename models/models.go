package models

import (
	"encoding/json"
	"time"
)

type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StripeID  string    `gorm:"uniqueIndex" json:"stripe_id"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaymentIntent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	StripeID    string    `gorm:"uniqueIndex" json:"stripe_id"`
	CustomerID  uint      `json:"customer_id"`
	Customer    Customer  `gorm:"foreignKey:CustomerID" json:"customer"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Metadata    JSONB     `gorm:"type:jsonb" json:"metadata"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Payment struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	StripeID        string        `gorm:"uniqueIndex" json:"stripe_id"`
	PaymentIntentID uint          `json:"payment_intent_id"`
	PaymentIntent   PaymentIntent `gorm:"foreignKey:PaymentIntentID" json:"payment_intent"`
	Amount          int64         `json:"amount"`
	Currency        string        `json:"currency"`
	Status          string        `json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
}

type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (interface{}, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return nil
	}

	if len(data) == 0 {
		*j = JSONB{}
		return nil
	}

	return json.Unmarshal(data, j)
}

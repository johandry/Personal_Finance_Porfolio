package models

import (
	"time"
)

// DebtType represents the type of debt
type DebtType string

const (
	DebtTypeCreditCard DebtType = "credit_card"
	DebtTypeLoan       DebtType = "loan"
	DebtTypeMortgage   DebtType = "mortgage"
	DebtTypeOther      DebtType = "other"
)

// Debt represents a financial debt
type Debt struct {
	ID           string    `json:"id"`
	Type         DebtType  `json:"type"`
	Name         string    `json:"name"`
	Principal    float64   `json:"principal"`
	CurrentValue float64   `json:"current_value"`
	Currency     string    `json:"currency"`
	InterestRate float64   `json:"interest_rate"`
	StartDate    time.Time `json:"start_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateDebtRequest represents the request body for creating a debt
type CreateDebtRequest struct {
	Type         DebtType `json:"type"`
	Name         string   `json:"name"`
	Principal    float64  `json:"principal"`
	CurrentValue *float64 `json:"current_value,omitempty"`
	Currency     string   `json:"currency"`
	InterestRate float64  `json:"interest_rate"`
	StartDate    string   `json:"start_date"`
}

// UpdateDebtRequest represents the request body for updating a debt
type UpdateDebtRequest struct {
	Name         *string  `json:"name,omitempty"`
	CurrentValue *float64 `json:"current_value,omitempty"`
	InterestRate *float64 `json:"interest_rate,omitempty"`
}

package models

import (
	"time"
)

// AssetType represents the type of asset
type AssetType string

const (
	AssetTypeStock      AssetType = "stock"
	AssetTypeProperty   AssetType = "property"
	AssetTypeCar        AssetType = "car"
	AssetTypeCash       AssetType = "cash"
	AssetTypeInvestment AssetType = "investment"
)

// AssetSource represents the data source
type AssetSource string

const (
	AssetSourceManual    AssetSource = "manual"
	AssetSourceMarketAPI AssetSource = "market_api"
)

// Asset represents a financial asset
type Asset struct {
	ID           string      `json:"id"`
	Type         AssetType   `json:"type"`
	Name         string      `json:"name"`
	BuyPrice     float64     `json:"buy_price"`
	CurrentValue float64     `json:"current_value"`
	Currency     string      `json:"currency"`
	Quantity     float64     `json:"quantity"`
	PurchaseDate time.Time   `json:"purchase_date"`
	Source       AssetSource `json:"source"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// AssetHistory represents historical values of an asset
type AssetHistory struct {
	ID        string    `json:"id"`
	AssetID   string    `json:"asset_id"`
	Value     float64   `json:"value"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

// ProfitLoss calculates the profit or loss for the asset
func (a *Asset) ProfitLoss() float64 {
	return (a.CurrentValue - a.BuyPrice) * a.Quantity
}

// TotalValue returns the total current value
func (a *Asset) TotalValue() float64 {
	return a.CurrentValue * a.Quantity
}

// CreateAssetRequest represents the request body for creating an asset
type CreateAssetRequest struct {
	Type         AssetType   `json:"type"`
	Name         string      `json:"name"`
	BuyPrice     float64     `json:"buy_price"`
	CurrentValue *float64    `json:"current_value,omitempty"`
	Currency     string      `json:"currency"`
	Quantity     float64     `json:"quantity"`
	PurchaseDate string      `json:"purchase_date"`
	Source       AssetSource `json:"source"`
}

// UpdateAssetRequest represents the request body for updating an asset
type UpdateAssetRequest struct {
	Name         *string      `json:"name,omitempty"`
	CurrentValue *float64     `json:"current_value,omitempty"`
	Quantity     *float64     `json:"quantity,omitempty"`
	Source       *AssetSource `json:"source,omitempty"`
}

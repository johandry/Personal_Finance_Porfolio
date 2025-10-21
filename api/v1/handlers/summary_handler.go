package handlers

import (
	"net/http"
	"time"

	"personal-finance/api/v1/db"
	"personal-finance/api/v1/models"
	"personal-finance/api/v1/services"
)

// SummaryHandler handles summary-related requests
type SummaryHandler struct {
	db         *db.PostgresDB
	marketData *services.MarketDataService
}

// NewSummaryHandler creates a new summary handler
func NewSummaryHandler(database *db.PostgresDB, marketDataService *services.MarketDataService) *SummaryHandler {
	return &SummaryHandler{
		db:         database,
		marketData: marketDataService,
	}
}

// GetNetWorth handles GET /api/v1/networth
func (h *SummaryHandler) GetNetWorth(w http.ResponseWriter, r *http.Request) {
	// Fetch all assets with real-time prices
	totalAssets := h.calculateTotalAssetsWithMarketData()

	// Calculate total debts
	var totalDebts float64
	debtQuery := `
		SELECT COALESCE(SUM(current_value), 0)
		FROM debts
	`
	err := h.db.DB.QueryRow(debtQuery).Scan(&totalDebts)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to calculate total debts")
		return
	}

	netWorth := models.NetWorth{
		TotalAssets:  totalAssets,
		TotalDebts:   totalDebts,
		NetWorth:     totalAssets - totalDebts,
		Currency:     "USD",
		CalculatedAt: time.Now(),
	}

	respondWithJSON(w, http.StatusOK, netWorth)
}

// calculateTotalAssetsWithMarketData fetches real-time prices for stocks and calculates total
func (h *SummaryHandler) calculateTotalAssetsWithMarketData() float64 {
	query := `
		SELECT id, type, name, buy_price, current_value, quantity, source
		FROM assets
	`

	rows, err := h.db.DB.Query(query)
	if err != nil {
		return 0
	}
	defer rows.Close()

	var total float64
	for rows.Next() {
		var id, assetType, name, source string
		var buyPrice, currentValue, quantity float64

		err := rows.Scan(&id, &assetType, &name, &buyPrice, &currentValue, &quantity, &source)
		if err != nil {
			continue
		}

		// Get real-time price for stocks
		if assetType == "stock" && source == "market_api" {
			price, err := h.marketData.GetCurrentValue(assetType, name, currentValue, source)
			if err == nil {
				currentValue = price
			}
		}

		total += currentValue * quantity
	}

	return total
}

// calculateTotalProfitLossWithMarketData calculates profit/loss with real-time prices
func (h *SummaryHandler) calculateTotalProfitLossWithMarketData() float64 {
	query := `
		SELECT type, name, buy_price, current_value, quantity, source
		FROM assets
	`

	rows, err := h.db.DB.Query(query)
	if err != nil {
		return 0
	}
	defer rows.Close()

	var totalProfitLoss float64
	for rows.Next() {
		var assetType, name, source string
		var buyPrice, currentValue, quantity float64

		err := rows.Scan(&assetType, &name, &buyPrice, &currentValue, &quantity, &source)
		if err != nil {
			continue
		}

		// Get real-time price for stocks
		if assetType == "stock" && source == "market_api" {
			price, err := h.marketData.GetCurrentValue(assetType, name, currentValue, source)
			if err == nil {
				currentValue = price
			}
		}

		totalProfitLoss += (currentValue - buyPrice) * quantity
	}

	return totalProfitLoss
}

// GetSummary handles GET /api/v1/summary
func (h *SummaryHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	// Calculate total assets with real-time prices
	totalAssets := h.calculateTotalAssetsWithMarketData()

	// Calculate total debts
	var totalDebts float64
	debtQuery := `
		SELECT COALESCE(SUM(current_value), 0)
		FROM debts
	`
	err := h.db.DB.QueryRow(debtQuery).Scan(&totalDebts)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to calculate total debts")
		return
	}

	// Calculate total profit/loss with real-time prices
	totalProfitLoss := h.calculateTotalProfitLossWithMarketData()

	// For daily profit/loss, we would need yesterday's data
	// For MVP, we'll set it to 0 or calculate from history if available
	var dailyProfitLoss float64 = 0.0

	summary := models.Summary{
		Date:            time.Now(),
		TotalAssets:     totalAssets,
		TotalDebts:      totalDebts,
		NetWorth:        totalAssets - totalDebts,
		DailyProfitLoss: dailyProfitLoss,
		TotalProfitLoss: totalProfitLoss,
		Currency:        "USD",
	}

	respondWithJSON(w, http.StatusOK, summary)
}

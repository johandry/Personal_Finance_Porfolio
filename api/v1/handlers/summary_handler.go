package handlers

import (
	"net/http"
	"time"

	"personal-finance/api/v1/db"
	"personal-finance/api/v1/models"
)

// SummaryHandler handles summary-related requests
type SummaryHandler struct {
	db *db.PostgresDB
}

// NewSummaryHandler creates a new summary handler
func NewSummaryHandler(database *db.PostgresDB) *SummaryHandler {
	return &SummaryHandler{db: database}
}

// GetNetWorth handles GET /api/v1/networth
func (h *SummaryHandler) GetNetWorth(w http.ResponseWriter, r *http.Request) {
	// Calculate total assets
	var totalAssets float64
	assetQuery := `
		SELECT COALESCE(SUM(current_value * quantity), 0)
		FROM assets
	`
	err := h.db.DB.QueryRow(assetQuery).Scan(&totalAssets)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to calculate total assets")
		return
	}

	// Calculate total debts
	var totalDebts float64
	debtQuery := `
		SELECT COALESCE(SUM(current_value), 0)
		FROM debts
	`
	err = h.db.DB.QueryRow(debtQuery).Scan(&totalDebts)
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

// GetSummary handles GET /api/v1/summary
func (h *SummaryHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	// Calculate total assets
	var totalAssets float64
	assetQuery := `
		SELECT COALESCE(SUM(current_value * quantity), 0)
		FROM assets
	`
	err := h.db.DB.QueryRow(assetQuery).Scan(&totalAssets)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to calculate total assets")
		return
	}

	// Calculate total debts
	var totalDebts float64
	debtQuery := `
		SELECT COALESCE(SUM(current_value), 0)
		FROM debts
	`
	err = h.db.DB.QueryRow(debtQuery).Scan(&totalDebts)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to calculate total debts")
		return
	}

	// Calculate total profit/loss
	var totalProfitLoss float64
	profitQuery := `
		SELECT COALESCE(SUM((current_value - buy_price) * quantity), 0)
		FROM assets
	`
	err = h.db.DB.QueryRow(profitQuery).Scan(&totalProfitLoss)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to calculate profit/loss")
		return
	}

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

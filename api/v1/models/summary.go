package models

import (
	"time"
)

// NetWorth represents the net worth summary
type NetWorth struct {
	TotalAssets  float64   `json:"total_assets"`
	TotalDebts   float64   `json:"total_debts"`
	NetWorth     float64   `json:"net_worth"`
	Currency     string    `json:"currency"`
	CalculatedAt time.Time `json:"calculated_at"`
}

// Summary represents daily summary of profit/loss and net worth
type Summary struct {
	Date            time.Time `json:"date"`
	TotalAssets     float64   `json:"total_assets"`
	TotalDebts      float64   `json:"total_debts"`
	NetWorth        float64   `json:"net_worth"`
	DailyProfitLoss float64   `json:"daily_profit_loss"`
	TotalProfitLoss float64   `json:"total_profit_loss"`
	Currency        string    `json:"currency"`
}

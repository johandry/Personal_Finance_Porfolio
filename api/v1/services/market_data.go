package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// MarketDataProvider defines the data source
type MarketDataProvider string

const (
	DefaultMarketDataProvider                    = ProviderYahooFinance
	ProviderYahooFinance      MarketDataProvider = "yahoo"
	ProviderAlphaVantage      MarketDataProvider = "alphavantage"
)

// MarketDataService provides stock market data
type MarketDataService struct {
	provider   MarketDataProvider
	apiKey     string
	httpClient *http.Client
	cache      map[string]*CachedPrice
	db         *sql.DB
}

// CachedPrice stores a price with timestamp
type CachedPrice struct {
	Price     float64
	Timestamp time.Time
}

// NewMarketDataService creates a new market data service
func NewMarketDataService(db *sql.DB) *MarketDataService {
	// Check which provider to use
	provider := os.Getenv("MARKET_DATA_PROVIDER")
	if provider == "" {
		provider = string(DefaultMarketDataProvider)
	}

	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		apiKey = "demo"
	}

	return &MarketDataService{
		provider: MarketDataProvider(provider),
		apiKey:   apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: make(map[string]*CachedPrice),
		db:    db,
	}
}

// GetStockPrice fetches the current price for a stock symbol
// It first checks the database cache, and only fetches from API if cache is older than 1 hour
func (s *MarketDataService) GetStockPrice(symbol string) (float64, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	// Check database cache first (valid for 1 hour)
	var price float64
	var lastUpdated time.Time
	checkQuery := `SELECT price, last_updated FROM stock_prices WHERE symbol = $1`
	err := s.db.QueryRow(checkQuery, symbol).Scan(&price, &lastUpdated)

	if err == nil {
		// Found in cache, check if still valid
		if time.Since(lastUpdated) < 60*time.Minute {
			fmt.Printf("[MarketData] Using DB cached price for %s: %.2f (age: %v)\n", symbol, price, time.Since(lastUpdated))
			return price, nil
		}
		fmt.Printf("[MarketData] DB cache expired for %s (age: %v), fetching fresh data\n", symbol, time.Since(lastUpdated))
	} else if err != sql.ErrNoRows {
		// Database error (not just "no rows"), log it but continue
		fmt.Printf("[MarketData] DB cache check error for %s: %v\n", symbol, err)
	}

	// Fetch from API
	var fetchErr error
	switch s.provider {
	case ProviderYahooFinance:
		price, fetchErr = s.getYahooFinancePrice(symbol)
	case ProviderAlphaVantage:
		price, fetchErr = s.getAlphaVantagePrice(symbol)
	default:
		price, fetchErr = s.getYahooFinancePrice(symbol)
	}

	if fetchErr != nil {
		return 0, fetchErr
	}

	// Store in database cache
	upsertQuery := `
		INSERT INTO stock_prices (symbol, price, last_updated, created_at)
		VALUES ($1, $2, $3, $3)
		ON CONFLICT (symbol) 
		DO UPDATE SET price = $2, last_updated = $3
	`
	now := time.Now()
	_, err = s.db.Exec(upsertQuery, symbol, price, now)
	if err != nil {
		fmt.Printf("[MarketData] Failed to cache price in DB for %s: %v\n", symbol, err)
		// Don't fail the request if caching fails
	} else {
		fmt.Printf("[MarketData] Cached %s price in DB: %.2f\n", symbol, price)
	}

	return price, nil
}

// getYahooFinancePrice fetches price from Yahoo Finance (FREE, no API key)
func (s *MarketDataService) getYahooFinancePrice(symbol string) (float64, error) {
	// Using Yahoo Finance query API (free, no authentication)
	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=1d",
		symbol,
	)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch stock price: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse Yahoo Finance response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	// Navigate the nested structure
	chart, ok := result["chart"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response format: missing chart")
	}

	resultArray, ok := chart["result"].([]interface{})
	if !ok || len(resultArray) == 0 {
		return 0, fmt.Errorf("invalid response format: missing result array")
	}

	firstResult, ok := resultArray[0].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response format: invalid result")
	}

	meta, ok := firstResult["meta"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response format: missing meta")
	}

	// Get current price from meta
	priceInterface, ok := meta["regularMarketPrice"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid response format: missing price")
	}

	return priceInterface, nil
}

// getAlphaVantagePrice fetches price from Alpha Vantage (requires API key)
func (s *MarketDataService) getAlphaVantagePrice(symbol string) (float64, error) {
	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		symbol, s.apiKey,
	)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch stock price: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for error messages
	if note, ok := result["Note"].(string); ok {
		return 0, fmt.Errorf("API limit reached: %s", note)
	}

	if errMsg, ok := result["Error Message"].(string); ok {
		return 0, fmt.Errorf("API error: %s", errMsg)
	}

	// Extract price from Global Quote
	globalQuote, ok := result["Global Quote"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response format: missing Global Quote")
	}

	priceStr, ok := globalQuote["05. price"].(string)
	if !ok {
		return 0, fmt.Errorf("invalid response format: missing price field")
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	return price, nil
}

// IsStockSymbol checks if an asset name is likely a stock symbol
func IsStockSymbol(name string) bool {
	// Simple heuristic: stock symbols are typically 1-5 uppercase letters
	name = strings.TrimSpace(name)
	if len(name) < 1 || len(name) > 5 {
		return false
	}

	for _, char := range name {
		if char < 'A' || char > 'Z' {
			return false
		}
	}

	return true
}

// GetCurrentValue returns the current value for an asset
// If it's a stock and market_api source, fetch from API
// Otherwise return the stored value
func (s *MarketDataService) GetCurrentValue(assetType, name string, storedValue float64, source string) (float64, error) {
	// Only fetch for stocks with market_api source
	if assetType == "stock" && source == "market_api" && IsStockSymbol(name) {
		fmt.Printf("[MarketData] Fetching real-time price for %s (stored: %.2f)\n", name, storedValue)
		price, err := s.GetStockPrice(name)
		if err != nil {
			// If API fails, fall back to stored value
			fmt.Printf("[MarketData] ERROR fetching %s: %v - using stored value\n", name, err)
			return storedValue, nil
		}
		fmt.Printf("[MarketData] Successfully fetched %s: %.2f\n", name, price)

		// Update all assets with this stock symbol in the database
		updateQuery := `
			UPDATE assets 
			SET current_value = $1, updated_at = $2 
			WHERE name = $3 AND type = 'stock' AND source = 'market_api'
		`
		result, err := s.db.Exec(updateQuery, price, time.Now(), name)
		if err != nil {
			fmt.Printf("[MarketData] Failed to update assets for %s: %v\n", name, err)
		} else {
			rowsAffected, _ := result.RowsAffected()
			fmt.Printf("[MarketData] Updated %d asset(s) with %s price\n", rowsAffected, name)
		}

		return price, nil
	}

	// For non-stocks or manual source, return stored value
	return storedValue, nil
}

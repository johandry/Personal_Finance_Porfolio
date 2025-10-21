package services

import (
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
}

// CachedPrice stores a price with timestamp
type CachedPrice struct {
	Price     float64
	Timestamp time.Time
}

// NewMarketDataService creates a new market data service
func NewMarketDataService() *MarketDataService {
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
	}
}

// GetStockPrice fetches the current price for a stock symbol
func (s *MarketDataService) GetStockPrice(symbol string) (float64, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	// Check cache (valid for 5 minutes)
	if cached, ok := s.cache[symbol]; ok {
		if time.Since(cached.Timestamp) < 5*time.Minute {
			return cached.Price, nil
		}
	}

	var price float64
	var err error

	// Fetch based on provider
	switch s.provider {
	case ProviderYahooFinance:
		price, err = s.getYahooFinancePrice(symbol)
	case ProviderAlphaVantage:
		price, err = s.getAlphaVantagePrice(symbol)
	default:
		price, err = s.getYahooFinancePrice(symbol)
	}

	if err != nil {
		return 0, err
	}

	// Cache the result
	s.cache[symbol] = &CachedPrice{
		Price:     price,
		Timestamp: time.Now(),
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
		price, err := s.GetStockPrice(name)
		if err != nil {
			// If API fails, fall back to stored value
			return storedValue, nil
		}
		return price, nil
	}

	// For non-stocks or manual source, return stored value
	return storedValue, nil
}

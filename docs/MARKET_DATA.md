# Stock Market Integration

This application integrates with multiple market data providers to fetch real-time stock prices for your portfolio.

## Features

- **Real-Time Stock Prices**: Automatically fetches current market prices for stocks
- **Multiple Providers**: Yahoo Finance (free, no key) or Alpha Vantage (requires key)
- **Smart Caching**: Caches prices for 5 minutes to reduce API calls
- **Graceful Fallback**: Uses stored values if API is unavailable
- **Asset Type Detection**: Automatically identifies stock symbols (1-5 uppercase letters)

## Providers

### Yahoo Finance (Default - Recommended for Development)

✅ **Completely FREE**  
✅ **No API key required**  
✅ **No rate limits**  
✅ **Real-time data**  
✅ **Works immediately**

**Best for**: Development, personal use, small portfolios

### Alpha Vantage (Alternative)

⚠️ **Requires API key**  
⚠️ **Rate limited** (5 requests/min, 500/day on free tier)  
✅ **More reliable for production**  
✅ **Official API with support**

**Best for**: Production deployments, large portfolios (with paid plan)

## Setup

### Option 1: Yahoo Finance (Default - No Setup Required!)

The application uses Yahoo Finance by default. **No configuration needed!**

Just start the application:

```bash
make up
```

### Option 2: Alpha Vantage (If You Need It)

1. **Get Your Free API Key**

   Visit [Alpha Vantage](https://www.alphavantage.co/support/#api-key) to get a free API key.

2. **Configure Provider and API Key**

   Edit the `.env` file:

   ```bash
   MARKET_DATA_PROVIDER=alphavantage
   ALPHA_VANTAGE_API_KEY=your_actual_api_key_here
   ```

3. **Restart the Application**

   ```bash
   make restart
   ```

## Usage

### Creating Stock Assets

When creating a stock asset:

1. **Set Type to "Stock"**
2. **Set Source to "Market API"**
3. **Use the stock symbol as the Name** (e.g., AAPL, TSLA, MSFT)
4. **Do NOT enter Current Value** - it will be fetched automatically

**Example:**

- Name: `AAPL`
- Type: `Stock`
- Buy Price: `150.00`
- Quantity: `10`
- Source: `Market API`
- Current Value: (leave empty - fetched from API)

### How It Works

1. **API Fetches**: When you view assets, the system fetches current prices from Alpha Vantage
2. **Caching**: Prices are cached for 5 minutes to avoid hitting rate limits
3. **Calculations**: Net worth, profit/loss, and totals use real-time prices
4. **Fallback**: If API fails, stored values are used

### Supported Asset Sources

- **Manual**: Manually enter and update current values
- **Market API**: Automatically fetch real-time prices (stocks only)

## API Limitations

### Yahoo Finance

✅ **No rate limits**  
✅ **No API key needed**  
✅ **Perfect for development**

**Note**: Yahoo Finance is an unofficial API. For production with guaranteed uptime, consider Alpha Vantage.

### Alpha Vantage (If Using)

The free tier has strict limits:

- **5 requests/minute**: Don't refresh too frequently
- **500 requests/day**: ~20 stock portfolio checks per day

**Best Practices**:

1. **Use Manual Source for Non-Stocks**: Properties, cars, cash, etc.
2. **Limit Stock Count**: Free tier works best with <10 stocks
3. **Don't Spam Refresh**: Prices are cached for 5 minutes
4. **Valid Symbols Only**: Use official ticker symbols (e.g., AAPL, not Apple)

**Upgrade Options**:

For larger portfolios or more frequent updates:

- **Premium Plan**: $49.99/month for unlimited API calls
- Visit [Alpha Vantage Pricing](https://www.alphavantage.co/premium/)

## Technical Details

### Market Data Service

Located at `api/v1/services/market_data.go`:

- **GetStockPrice()**: Fetches price from Alpha Vantage API
- **GetCurrentValue()**: Returns real-time or stored value
- **IsStockSymbol()**: Validates stock symbol format
- **Caching**: In-memory cache with 5-minute TTL

### API Integration

The service is integrated into:

- `AssetHandler.ListAssets()` - Real-time prices in asset list
- `AssetHandler.GetAsset()` - Real-time price for single asset
- `SummaryHandler.GetNetWorth()` - Real-time net worth calculation
- `SummaryHandler.GetSummary()` - Real-time profit/loss calculation

### Database Schema

No changes to database schema required:

- `current_value` is still stored (as fallback)
- Stock symbols stored in `name` field
- `source` field determines whether to fetch from API

## Troubleshooting

### "API limit reached" Error

**Cause**: Exceeded 5 requests/minute or 500 requests/day

**Solution**:

- Wait 1 minute (for per-minute limit)
- Wait until next day (for daily limit)
- Upgrade to Premium plan

### Invalid Stock Symbol

**Cause**: Symbol not recognized by Alpha Vantage

**Solution**:

- Verify the symbol on [Yahoo Finance](https://finance.yahoo.com/)
- Use the correct exchange symbol
- Example: Use `TSLA` not `Tesla`

### "Failed to fetch stock price"

**Cause**: Network issues or API unavailable

**Solution**:

- Check internet connection
- Verify API key is correct
- System will use stored values as fallback

### Prices Not Updating

**Cause**: Cache is active (5-minute TTL)

**Solution**:

- Wait 5 minutes for cache to expire
- Prices update automatically after cache expiry
- This is normal behavior to avoid rate limits

## Alternative Market Data Providers

If you need a different provider, you can modify `api/v1/services/market_data.go`:

### Other Free APIs

- **Yahoo Finance** ✅ (already integrated, default)
- **IEX Cloud** (free tier: 50,000 messages/month)
- **Finnhub** (free tier: 60 calls/minute)
- **Twelve Data** (free tier: 800 requests/day)

### Implementation

```go
// Replace GetStockPrice() implementation with your preferred API
func (s *MarketDataService) GetStockPrice(symbol string) (float64, error) {
    // Your custom API integration here
}
```

## Future Enhancements

Potential improvements for Phase 2:

- [ ] Support for multiple market data providers
- [ ] Automatic daily price updates via cron job
- [ ] Historical price charts
- [ ] Currency conversion for international stocks
- [ ] Crypto asset support (Coinbase, Binance APIs)
- [ ] Real-time WebSocket updates
- [ ] Portfolio performance analytics

## Resources

- [Alpha Vantage Documentation](https://www.alphavantage.co/documentation/)
- [API Key Registration](https://www.alphavantage.co/support/#api-key)
- [Supported Stock Symbols](https://www.alphavantage.co/query?function=LISTING_STATUS&apikey=demo)
- [Rate Limit Details](https://www.alphavantage.co/premium/)

---

**Note**: The `demo` API key is for testing only and has very limited functionality. Get your free API key to use the full features!

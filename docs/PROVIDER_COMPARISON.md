# Market Data Provider Comparison

Quick comparison of available stock market data providers for development and production use.

## Quick Recommendation

- **Development**: Use Yahoo Finance (default, no setup)
- **Production**: Use Yahoo Finance or upgrade to Alpha Vantage Premium

## Provider Comparison

| Feature | Yahoo Finance | Alpha Vantage (Free) | Alpha Vantage (Premium) |
|---------|--------------|---------------------|------------------------|
| **Cost** | FREE | FREE | $49.99/month |
| **API Key** | ❌ Not required | ✅ Required | ✅ Required |
| **Rate Limit** | None | 5/min, 500/day | Unlimited |
| **Real-time Data** | ✅ Yes | ✅ Yes | ✅ Yes |
| **Setup Time** | 0 seconds | ~2 minutes | ~2 minutes |
| **Reliability** | Good (unofficial) | Good | Excellent |
| **Support** | Community | Limited | Priority |
| **Best For** | Development, personal | Small projects | Production, large apps |

## When to Use Each Provider

### Yahoo Finance (Default)

**Use when**:

- 🚀 Developing locally
- 💰 Need completely free solution
- 🏠 Personal finance tracking
- 📊 Small portfolio (<50 stocks)
- ⚡ Want zero setup time

**Don't use when**:

- 🏢 Mission-critical production app
- 📞 Need guaranteed support
- 📜 Require SLA/uptime guarantees

### Alpha Vantage Free

**Use when**:

- 🔒 Need official API with support
- 📋 Require terms of service
- 🏢 Small production deployment
- 📊 Portfolio with <10 stocks
- ✅ Can work within rate limits

**Don't use when**:

- ⚡ Need frequent updates (>5/min)
- 📈 Large portfolio (>10 stocks)
- 🔄 Need real-time streaming

### Alpha Vantage Premium

**Use when**:

- 🏢 Production application
- 📊 Large portfolio (>10 stocks)
- ⚡ Need frequent updates
- 🔒 Require guaranteed uptime
- 💼 Commercial application

## Switching Providers

Simply change the environment variable in `.env`:

```bash
# Use Yahoo Finance (default, no API key needed)
MARKET_DATA_PROVIDER=yahoo

# Use Alpha Vantage (requires API key)
MARKET_DATA_PROVIDER=alphavantage
ALPHA_VANTAGE_API_KEY=your_actual_key_here
```

Then restart:

```bash
make restart
```

## Testing Each Provider

### Test Yahoo Finance

```bash
# Set provider
echo "MARKET_DATA_PROVIDER=yahoo" >> .env

# Start app
make up

# Test stock price fetch
curl http://localhost:8080/api/v1/assets
```

### Test Alpha Vantage

```bash
# Set provider and key
echo "MARKET_DATA_PROVIDER=alphavantage" >> .env
echo "ALPHA_VANTAGE_API_KEY=your_key" >> .env

# Start app
make restart

# Test stock price fetch
curl http://localhost:8080/api/v1/assets
```

## Performance Comparison

Based on typical usage (checking 5 stocks):

| Provider | Response Time | Cache Hit | Cache Miss |
|----------|--------------|-----------|------------|
| Yahoo Finance | ~200ms | ~1ms | ~200ms |
| Alpha Vantage | ~300ms | ~1ms | ~300ms |

**Note**: Both providers cache results for 5 minutes, so repeated requests are instant.

## Supported Stock Symbols

Both providers support:

- ✅ US stocks (NYSE, NASDAQ)
- ✅ Major international exchanges
- ✅ ETFs
- ✅ Mutual funds

**Example symbols**:

- `AAPL` - Apple Inc.
- `TSLA` - Tesla, Inc.
- `MSFT` - Microsoft Corporation
- `GOOGL` - Alphabet Inc.
- `SPY` - S&P 500 ETF

## Data Accuracy

Both providers offer:

- Real-time prices (15-20 minute delay for free tiers)
- Accurate market data
- Regular updates during market hours

**Market Hours**: 9:30 AM - 4:00 PM ET (NYSE/NASDAQ)

## Error Handling

Both providers include:

- Automatic fallback to stored values
- 5-minute caching to reduce API calls
- Graceful error handling
- Clear error messages

## Cost Analysis

### Development (1 developer)

- **Yahoo Finance**: $0/month ✅
- **Alpha Vantage Free**: $0/month
- **Alpha Vantage Premium**: $49.99/month

### Small Team (3-5 developers)

- **Yahoo Finance**: $0/month ✅
- **Alpha Vantage Free**: Limited (may hit rate limits)
- **Alpha Vantage Premium**: $49.99/month

### Production (100+ users)

- **Yahoo Finance**: $0/month (monitor for issues)
- **Alpha Vantage Free**: Not recommended (rate limits)
- **Alpha Vantage Premium**: $49.99/month ✅

## Troubleshooting

### Yahoo Finance Issues

If Yahoo Finance fails:

1. Check internet connection
2. Verify stock symbol exists
3. Try Alpha Vantage as backup

```bash
MARKET_DATA_PROVIDER=alphavantage
```

### Alpha Vantage Issues

If hitting rate limits:

1. Switch to Yahoo Finance temporarily
2. Increase cache duration
3. Upgrade to Premium plan

## Future Providers

Potential additions in future releases:

- [ ] Finnhub (free tier: 60 calls/min)
- [ ] IEX Cloud (free tier: 50k messages/month)
- [ ] Twelve Data (free tier: 800 requests/day)
- [ ] Polygon.io (paid, unlimited)

## Conclusion

**For most users**: Stick with Yahoo Finance (default). It's free, fast, and requires zero setup.

**For production apps**: Consider Alpha Vantage Premium for guaranteed uptime and support.

---

**Current Default**: Yahoo Finance  
**Recommended**: Yahoo Finance for development, Alpha Vantage Premium for production

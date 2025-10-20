# Quick Start Guide

Get your Personal Finance Portfolio up and running in minutes!

## 🚀 Option 1: Docker (Recommended)

### Prerequisites

- Docker Desktop installed
- Docker Compose installed

### Steps

1. **Start everything**

   ```bash
   make up
   ```

2. **Open in browser**
   - Frontend: <http://localhost:3000>
   - API: <http://localhost:8080>

3. **Start tracking your finances!**
   - Go to <http://localhost:3000>
   - Click "Manage Assets & Debts"
   - Add your first asset

That's it! 🎉

## 📱 Using the Application

### Adding Your First Asset

1. Navigate to **Manage Assets & Debts** page
2. Click **"+ Add New Asset"**
3. Fill in the details:

   ```text
   Name: Apple Stock
   Type: Stock
   Buy Price: 150.00
   Current Value: 175.50
   Quantity: 10
   Currency: USD
   Purchase Date: 2024-01-01
   Source: Manual
   ```

4. Click **"Save Asset"**

You'll see:

- Total value: $1,755.00
- Profit: $255.00 (green ✅)

### Viewing Your Dashboard

The dashboard automatically shows:

- 📊 Total Assets
- 📉 Total Debts
- 💎 Net Worth
- 📈 Profit/Loss

Plus beautiful charts showing your asset distribution!

## 🎯 Common Tasks

### Update Asset Value

1. Go to **Manage** page
2. Find your asset
3. Click **"Edit"**
4. Update **Current Value**
5. Click **"Save"**

The dashboard updates automatically!

### Add a Debt

1. Go to **Manage** page
2. Click **Debts** tab
3. Click **"+ Add New Debt"**
4. Fill in details:

   ```text
   Name: Car Loan
   Type: Loan
   Principal: 25,000.00
   Current Value: 20,000.00
   Interest Rate: 4.5
   Start Date: 2023-06-01
   ```

5. Click **"Save Debt"**

### Delete Items

1. Find the item in the table
2. Click **"Delete"**
3. Confirm deletion

## 🛠️ Troubleshooting

### Port Already in Use

If you get a port error:

```bash
# Stop existing containers
make down

# Or change ports in docker-compose.yml
ports:
  - "3001:80"  # Change 3000 to 3001
```

### Can't Connect to API

1. Check API is running:

   ```bash
   make logs-api
   ```

2. Test API directly:

   ```bash
   curl http://localhost:8080/api/v1/health
   ```

3. Restart services:

   ```bash
   make restart
   ```

### Frontend Not Loading

1. Check web container:

   ```bash
   make logs-web
   ```

2. Clear browser cache
3. Try incognito mode

## 📊 Understanding Your Data

### Profit/Loss Calculation

```text
Profit/Loss = (Current Value - Buy Price) × Quantity
```

Example:

- Buy Price: $100
- Current Value: $120
- Quantity: 10
- **Profit: $200** ✅

### Net Worth

```text
Net Worth = Total Assets - Total Debts
```

Example:

- Assets: $50,000
- Debts: $15,000
- **Net Worth: $35,000** 💎

## 🔧 Advanced Usage

### Access the Database

```bash
make shell-db
```

Then run SQL:

```sql
SELECT * FROM assets;
SELECT * FROM debts;
```

### View API Logs

```bash
make logs-api
```

### Run Tests

```bash
make test
```

## 🎨 Customization

### Change API URL

Edit `web/js/config.js`:

```javascript
const API_CONFIG = {
    baseURL: 'http://your-api-url:8080/api/v1',
};
```

### Change Colors

Edit `web/css/styles.css`:

```css
:root {
    --primary-color: #your-color;
}
```

## 📱 Mobile Usage

The application is fully responsive! Just:

1. Open <http://localhost:3000> on your phone
2. Add to home screen for app-like experience
3. Track finances on the go!

## 🔒 Security Tips

For production deployment:

1. Add authentication to the API
2. Use HTTPS
3. Set strong database passwords
4. Enable firewall rules
5. Regular backups

## 💡 Tips & Tricks

### Asset Entry Tips

- Start with major assets (house, car, investments)
- Update stock values weekly
- Set realistic current values for property
- Track cash in different accounts separately

### Dashboard Usage

- Check dashboard daily for overview
- Watch the profit/loss trend
- Review asset distribution monthly
- Adjust investments based on allocation

### Performance

- The dashboard auto-refreshes every 30 seconds
- Historical data is stored automatically
- Charts update in real-time

## 🆘 Getting Help

1. Check the main [README](../README.md)
2. Check browser console for errors
3. View container logs: `make logs`

## 🎉 Next Steps

1. ✅ Add all your assets
2. ✅ Add all your debts
3. ✅ Review your net worth
4. ✅ Set financial goals
5. ✅ Track progress monthly
6. ✅ Make informed decisions!

---

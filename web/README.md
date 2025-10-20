# Web Frontend

This directory contains the static frontend for the Personal Finance Portfolio application.

## 📁 Structure

```tree
web/
├── index.html          # Dashboard page with charts and summaries
├── manage.html         # CRUD interface for assets and debts
├── css/
│   └── styles.css      # All styling for the application
└── js/
    ├── config.js       # Configuration and utility functions
    ├── api.js          # API client for backend communication
    ├── dashboard.js    # Dashboard page logic
    └── manage.js       # Management page logic
```

## 🎨 Pages

### Dashboard (`index.html`)

The main landing page featuring:

- **Summary Cards:** Total assets, debts, net worth, and profit/loss
- **Charts:**
  - Asset distribution pie chart
  - Net worth bar chart
- **Recent Items:** Quick view of recent assets and debts
- Real-time updates every 30 seconds

### Manage (`manage.html`)

Full CRUD interface with:

- **Asset Management:** Create, read, update, delete assets
- **Debt Management:** Create, read, update, delete debts
- **Modal Forms:** User-friendly forms for data entry
- **Validation:** Client-side form validation
- **Toast Notifications:** Success and error feedback

## 🚀 Features

### UI/UX

- ✨ Modern, clean design
- 📱 Fully responsive (mobile, tablet, desktop)
- 🎨 Professional color scheme
- 🌊 Smooth animations and transitions
- 💬 Toast notifications for user feedback

### Functionality

- 📊 Interactive charts with Chart.js
- 🔄 Real-time data updates
- 💰 Currency formatting
- 📅 Date formatting
- 🏷️ Type badges with color coding
- ➕➖ Profit/loss color indicators

### Asset Types Supported

- 📈 Stocks
- 🏠 Property
- 🚗 Cars
- 💵 Cash
- 📊 Investments

### Debt Types Supported

- 💳 Credit Cards
- 🏦 Loans
- 🏡 Mortgages
- 📋 Other

## 🔧 Configuration

The API endpoint is configured in `js/config.js`:

```javascript
const API_CONFIG = {
    baseURL: 'http://localhost:8080/api/v1',
    timeout: 10000,
};
```

To change the API endpoint, edit this file before deployment.

## 🌐 Deployment

### With Docker (Recommended)

The frontend is automatically served via Nginx when using docker-compose:

```bash
docker-compose up
```

Access at: <http://localhost:3000>

### Standalone Static Hosting

You can host the `web/` directory on any static file server:

#### Python

```bash
cd web
python3 -m http.server 3000
```

#### Node.js (with http-server)

```bash
npm install -g http-server
cd web
http-server -p 3000
```

#### Nginx

```nginx
server {
    listen 80;
    root /path/to/web;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

## 🎯 Usage Guide

### Adding an Asset

1. Navigate to the **Manage** page
2. Click **"+ Add New Asset"**
3. Fill in the form:
   - Name (e.g., "Tesla Stock")
   - Type (Stock, Property, Car, etc.)
   - Buy Price
   - Current Value (optional, defaults to buy price)
   - Quantity
   - Currency
   - Purchase Date
   - Source (Manual or Market API)
4. Click **"Save Asset"**

### Updating an Asset

1. Go to the **Manage** page
2. Find the asset in the table
3. Click **"Edit"**
4. Modify the values
5. Click **"Save Asset"**

### Deleting an Asset

1. Navigate to the **Manage** page
2. Find the asset in the table
3. Click **"Delete"**
4. Confirm the deletion

The same process applies to debts.

## 📊 Dashboard Insights

### Summary Cards

- **Total Assets:** Sum of all asset values
- **Total Debts:** Sum of all debt values
- **Net Worth:** Assets minus debts
- **Total Profit/Loss:** Cumulative gains/losses

### Charts

- **Asset Distribution:** Pie chart showing allocation by type
- **Net Worth Overview:** Bar chart comparing invested vs current value

## 🎨 Customization

### Changing Colors

Edit `css/styles.css` and modify the CSS variables:

```css
:root {
    --primary-color: #4f46e5;      /* Main brand color */
    --success-color: #10b981;      /* Positive values */
    --danger-color: #ef4444;       /* Negative values */
    --background: #f8fafc;         /* Page background */
    /* ... more variables */
}
```

### Adding New Asset Types

1. Update the dropdown in `manage.html`:

    ```html
    <select id="assetType" required>
        <option value="your_new_type">Your New Type</option>
    </select>
    ```

2. Add styling in `css/styles.css`:

    ```css
    .badge.your_new_type { 
        background: #color; 
        color: #textcolor; 
    }
    ```

## 🔐 Security Notes

- This is a client-side only application
- All data is fetched from the backend API
- No authentication is implemented in MVP
- For production, add authentication headers in `api.js`

## 🐛 Troubleshooting

### Charts Not Displaying

- Check browser console for errors
- Ensure Chart.js is loaded from CDN
- Verify API is returning data

### API Connection Failed

- Confirm API is running on port 8080
- Check CORS settings in backend
- Verify `API_CONFIG.baseURL` in `config.js`

### Data Not Loading

- Open browser DevTools Network tab
- Check for failed API requests
- Verify backend is accessible
- Check browser console for JavaScript errors

## 📱 Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)
- Mobile browsers (iOS Safari, Chrome Mobile)

## 🔮 Future Enhancements

- [ ] Dark mode support
- [ ] Export data to CSV/PDF
- [ ] Advanced filtering and search
- [ ] Historical trend charts
- [ ] Budget tracking
- [ ] Multi-currency support
- [ ] Offline support with Service Workers
- [ ] Real-time market data integration

## 📄 License

MIT License - see main README

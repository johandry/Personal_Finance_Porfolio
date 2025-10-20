# 🧩 Personal Finance Portfolio

> A simple API-first system for tracking, visualizing, and updating financial assets and debts with profit/loss calculations over time.

## 📋 Features

- **Asset Management:** Track stocks, properties, cars, cash, and investments
- **Debt Tracking:** Monitor credit cards, loans, mortgages, and other debts
- **Historical Tracking:** Store and view daily asset values
- **Net Worth Calculation:** Automatic aggregation of assets minus debts
- **Profit/Loss Analysis:** Daily delta and cumulative returns
- **Interactive Dashboard:** Beautiful charts and visualizations
- **CRUD Interface:** Easy-to-use web interface for managing finances
- **RESTful API:** Built with Go for high performance
- **Containerized:** Docker & Docker Compose for easy deployment

## 🛠️ Tech Stack

### Backend

- **Language:** Go
- **Framework:** Chi Router
- **Database:** PostgreSQL

### Frontend

- **Static Site:** HTML5, CSS3, JavaScript (ES6+)
- **Charts:** Chart.js
- **Web Server:** Nginx

### DevOps

- **Containerization:** Docker & Docker Compose
- **API Documentation:** OpenAPI (Swagger) ready

## 📁 Project Structure

```tree
.
├── api/                         # Backend API
│   └── v1/
│       ├── handlers/            # HTTP request handlers
│       │   ├── asset_handler.go
│       │   ├── debt_handler.go
│       │   ├── summary_handler.go
│       │   └── utils.go
│       ├── models/              # Data models
│       │   ├── asset.go
│       │   ├── debt.go
│       │   └── summary.go
│       └── db/                  # Database layer
│           └── postgres.go
├── web/                         # Frontend
│   ├── index.html               # Dashboard page
│   ├── manage.html              # CRUD page
│   ├── css/
│   │   └── styles.css           # Styles
│   └── js/
│       ├── config.js            # Configuration
│       ├── api.js               # API client
│       ├── dashboard.js         # Dashboard logic
│       └── manage.js            # CRUD logic
├── docs/
│   └── prd.md                   # Product Requirements Document
├── main.go                      # Application entry point
├── go.mod                       # Go dependencies
├── go.sum                       # Dependency checksums
├── .env                         # Environment variables template
├── Dockerfile                   # Docker configuration
├── docker-compose.yml           # Multi-container Docker config
├── nginx.conf                   # Nginx configuration
├── Makefile                     # Build and run automation
└── README.md                    # This file
```

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Make (optional, for convenience)
- Go   (for local development)

### 1. Clone the Repository

```bash
git clone git@github.com:johandry/Personal_Finance_Porfolio.git
cd Personal_Finance_Porfolio
```

### 2. Configure Environment

Edit api/.env if needed (defaults work with Docker Compose)

### 3. Start the Application

Using Make:

```bash
make up
```

Or using Docker Compose directly:

```bash
docker-compose up --build
```

### 4. Access the Application

| Service  | URL                           | Description              |
|----------|-------------------------------|--------------------------|
| Frontend | <http://localhost:3000>         | Web interface            |
| API      | <http://localhost:8080>         | REST API                 |
| Database | localhost:5432                | PostgreSQL               |
| Health   | <http://localhost:8080/api/v1/health> | API health check |

### 5. Test the System

Open your browser and navigate to:

- **Dashboard:** <http://localhost:3000/index.html>
- **Manage Assets:** <http://localhost:3000/manage.html>

Or test the API directly:

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Create an asset
curl -X POST http://localhost:8080/api/v1/assets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "stock",
    "name": "AAPL",
    "quantity": 20,
    "buy_price": 150.0,
    "current_value": 172.35,
    "currency": "USD",
    "purchase_date": "2024-07-01",
    "source": "market_api"
  }' | jq
```

## 📚 API Endpoints

### Assets

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/api/v1/assets` | Create new asset |
| GET    | `/api/v1/assets` | List all assets |
| GET    | `/api/v1/assets/{id}` | Get specific asset |
| PUT    | `/api/v1/assets/{id}` | Update asset |
| DELETE | `/api/v1/assets/{id}` | Delete asset |
| GET    | `/api/v1/assets/{id}/history` | Get asset history |

### Debts

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/api/v1/debts` | Create new debt |
| GET    | `/api/v1/debts` | List all debts |
| GET    | `/api/v1/debts/{id}` | Get specific debt |
| PUT    | `/api/v1/debts/{id}` | Update debt |
| DELETE | `/api/v1/debts/{id}` | Delete debt |

### Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/networth` | Get total net worth |
| GET | `/api/v1/summary` | Get daily summary with P/L |

## 🔧 Development

### Local Development (without Docker)

1. Start PostgreSQL:

    ```bash
    docker run --name finance-postgres \
      -e POSTGRES_USER=financeuser \
      -e POSTGRES_PASSWORD=financepass \
      -e POSTGRES_DB=financedb \
      -p 5432:5432 \
      -d postgres:15-alpine
    ```

2. Run the application:

    ```bash
    go mod download
    DB_HOST=localhost go run main.go
    ```

## 🛠️ Available Commands

```bash
make help       # Show all commands
make up         # Start all services
make down       # Stop all services
make restart    # Restart services
make logs       # View all logs
make logs-api   # API logs only
make logs-web   # Web logs only
make logs-db    # Database logs only
make clean      # Clean everything
make test       # Run tests
make dev        # Run API locally
make dev-web    # Serve frontend locally
make open       # Open in browser
make health     # Check API health
make status     # Service status
```

## 📊 Database Schema

### Assets Table

- `id` (UUID, Primary Key)
- `type` (VARCHAR: stock, property, car, cash, investment)
- `name` (VARCHAR)
- `buy_price` (DECIMAL)
- `current_value` (DECIMAL)
- `currency` (VARCHAR)
- `quantity` (DECIMAL)
- `purchase_date` (DATE)
- `source` (VARCHAR: manual, market_api)
- `created_at`, `updated_at` (TIMESTAMP)

### Asset History Table

- `id` (UUID, Primary Key)
- `asset_id` (UUID, Foreign Key)
- `value` (DECIMAL)
- `date` (DATE)
- `created_at` (TIMESTAMP)

### Debts Table

- `id` (UUID, Primary Key)
- `type` (VARCHAR: credit_card, loan, mortgage, other)
- `name` (VARCHAR)
- `principal` (DECIMAL)
- `current_value` (DECIMAL)
- `currency` (VARCHAR)
- `interest_rate` (DECIMAL)
- `start_date` (DATE)
- `created_at`, `updated_at` (TIMESTAMP)

## 🎨 Design Highlights

### Colors

- Primary: Indigo (#4f46e5)
- Success: Green (#10b981)
- Danger: Red (#ef4444)
- Background: Light Gray (#f8fafc)

### Responsive

- Desktop: Full layout
- Tablet: Grid adjusts
- Mobile: Single column, touch-friendly

### Features

- Smooth animations
- Hover effects
- Toast notifications
- Modal dialogs
- Loading states
- Empty states
- Error handling

## 🔧 Technology Stack

| Layer       | Technology                    |
|-------------|-------------------------------|
| Frontend    | HTML5, CSS3, JavaScript ES6+  |
| Charts      | Chart.js                      |
| Backend     | Go 1.21+                      |
| Framework   | Chi Router                    |
| Database    | PostgreSQL 15                 |
| Web Server  | Nginx (Alpine)                |
| Container   | Docker & Docker Compose       |

## 📖 Documentation

1. **README.md** - Main documentation
2. **docs/prd.md** - Product requirements
3. **docs/QUICKSTART.md** - Quick start guide
4. **web/README.md** - Frontend documentation
5. **Makefile** - All available commands

## 🎯 Roadmap

### Phase 1 (MVP) ✅

- [x] RESTful API with Go
- [x] PostgreSQL database
- [x] Docker containerization
- [x] Asset & Debt CRUD
- [x] Net worth calculation
- [x] Historical tracking
- [x] Static web frontend (HTML/CSS/JS)

### Phase 2 (Future)

- [ ] Market data API integration (Yahoo Finance, Alpha Vantage)
- [ ] Automated daily valuation updates
- [ ] Enhanced data visualizations
- [ ] Export to CSV/PDF
- [ ] Mobile app (Flutter)
- [ ] Bank integrations (Plaid)
- [ ] Notifications & insights
- [ ] Budget tracking
- [ ] Bill reminders
- [ ] Investment recommendations

### Phase 3 (Further Future)

- [ ] AI insights & recommendations

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.

## 📞 Support

For questions or issues, please open an issue on GitHub.

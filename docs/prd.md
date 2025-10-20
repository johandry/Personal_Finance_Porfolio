
# 🧩 Personal Finance Portfolio — MVP PRD

> Technical Specification | Go API + Static Web Frontend

---

## 1. Overview

**Goal:**  
Build a simple system for users to **track, visualize, and update** their financial assets and debts, and **calculate profit/loss over time**.

**Deliverables:**

- RESTful API in **Go**, containerized with **Docker**.
- Static frontend in **HTML/CSS/JS** for:
  - Viewing summaries and charts
  - CRUD for assets and debts

**Future expansion:** Mobile app (Phase 2).

---

## 2. Objectives

- Centralize all personal finance data.  
- Daily valuation for each asset.  
- Profit/loss tracking by date.  
- API-first approach to enable external integrations.

---

## 3. Target Users

- Individuals tracking investments or assets manually.  
- Developers who want an API-driven finance tool.  
- Use cases:
  - Record stock purchases and see profit/loss daily.
  - Add property, vehicles, or cash manually.
  - Track net worth growth.

---

## 4. Core MVP Features

### 4.1 Asset Categories

- **Stocks:** auto-valued via market API  
- **Properties / Cars:** manual or external valuation  
- **Bank Accounts / Cash:** manual values  
- **Debts / Credit Cards:** manual entries  
- **Investments:** generic type (manual or API)

### 4.2 Core Functions

| Function | Description |
|-----------|-------------|
| CRUD Assets/Debts | Create, read, update, delete any record |
| Historical Tracking | Store and view daily values |
| Net Worth Summary | Aggregate assets - debts |
| Profit/Loss | Daily delta and cumulative returns |
| Data Sync | Background task to refresh external values |

---

## 5. API Specification

### 5.1 Stack

- **Language:** Go  
- **Framework:** `net/http` and `chi`  
- **Database:** PostgreSQL  
- **Containerization:** Docker  
- **Docs:** OpenAPI (Swagger)

### 5.2 Base URL

`/api/v1`

### 5.3 Endpoints

#### Assets

| Method | Endpoint | Description |
|--------|-----------|-------------|
| POST   | `/assets` | Create new asset |
| GET    | `/assets` | List all assets |
| GET    | `/assets/{id}` | Retrieve specific asset |
| PUT    | `/assets/{id}` | Update asset |
| DELETE | `/assets/{id}` | Delete asset |
| GET    | `/assets/{id}/history` | Get historical values |

#### Debts

| Method | Endpoint | Description |
|--------|-----------|-------------|
| POST   | `/debts` | Create new debt |
| GET    | `/debts` | List all debts |

#### Summary

| Method | Endpoint | Description |
|--------|-----------|-------------|
| GET | `/networth` | Total assets - debts |
| GET | `/summary` | Daily summary of P/L and net worth |

---

### 5.4 Data Model (JSON)

```json
{
  "id": "uuid",
  "type": "stock|property|car|cash|debt|investment",
  "name": "Tesla Stock",
  "buy_price": 250.0,
  "current_value": 270.5,
  "currency": "USD",
  "quantity": 10,
  "purchase_date": "2024-06-01",
  "source": "manual|market_api",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-02T00:00:00Z"
}
```

**Profit/Loss Formula:**

`(current_value - buy_price) * quantity`

---

## 6. Web Frontend (Static)

### 6.1 Stack

- Plain **HTML / CSS / JS**
- Fetch API calls → Go backend
- Deployed as static files (can be served by Go or CDN)

### 6.2 Pages

#### `/index.html`

- Dashboard view
  - Total assets, debts, net worth
  - Simple charts or tables

#### `/manage.html`

- CRUD UI for assets/debts
- Form submissions using Fetch API

---

## 7. Architecture & Patterns

- **Separation of Concerns:**  
    Go handles business logic & persistence.  
    Frontend is static and stateless.

- **API-First:**  
    Same API for website and future mobile app.

- **Data Update Strategy:**
  - Cron or background job refreshes external data daily.
  - Manual overrides allowed.

- **Authentication (MVP):**  
    Token-based (JWT or static API key).

- **Error Format:**
    `{ "error": "invalid asset id" }`

---

## 8. Example Usage

### 8.1 Create an Asset

```bash
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

**Response:**

`{   "id": "uuid",   "name": "AAPL",   "current_value": 172.35,   "profit_loss": 445.0 }`

---

## 9. Metrics of Success

- ✅ User can add and update assets/debts manually.
- ✅ API calculates profit/loss daily.
- ✅ API uptime ≥ 99%.
- ✅ Frontend loads within 1s and fetches all data correctly.

---

## 10. 🚀 Roadmap

## Phase 1 (MVP)

- [x] REST API (Go)
- [x] CRUD assets/debts
- [x] Manual asset values
- [x] Static web dashboard

## Phase 2

- [ ] Bank & brokerage integrations (Plaid, etc.)
- [ ] User authentication
- [ ] Mobile app
- [ ] Data visualization (charts, trends)
- [ ] Notifications & insights

## Phase 3

- [ ] AI insights & recommendations

---

## 11. References

- Market data: e.g., Yahoo Finance API, Alpha Vantage
- Real estate data: Zillow, Redfin (future)
- Vehicle data: Kelley Blue Book (future)

---

## 12. Developer Notes

- Use `docker-compose` for local dev (Go + Postgres).
- Include `.env` for secrets & API keys.
- Generate Swagger docs automatically from Go comments.
- Store daily valuations in `asset_history` table.

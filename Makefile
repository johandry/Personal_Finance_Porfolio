.PHONY: help up down build restart logs clean test dev

# Default target
help:
	@echo "Personal Finance Portfolio - Makefile Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make up         - Start all services with Docker Compose"
	@echo "  make down       - Stop all services"
	@echo "  make build      - Build the API binary"
	@echo "  make restart    - Restart all services"
	@echo "  make logs       - View logs from all services"
	@echo "  make logs-api   - View logs from API service only"
	@echo "  make logs-db    - View logs from database service only"
	@echo "  make logs-web   - View logs from web service only"
	@echo "  make clean      - Stop services and remove volumes"
	@echo "  make test       - Run tests"
	@echo "  make dev        - Run API locally (requires local PostgreSQL)"
	@echo "  make dev-web    - Serve frontend locally with Python"
	@echo "  make migrate    - Run database migrations"
	@echo "  make shell-api  - Open shell in API container"
	@echo "  make shell-db   - Open PostgreSQL shell"
	@echo "  make open       - Open frontend in browser"
	@echo ""

# Start all services
up:
	@echo "🚀 Starting services..."
	docker-compose up --build -d
	@echo "✅ Services started!"
	@echo "Frontend available at: http://localhost:3000"
	@echo "API available at: http://localhost:8080"
	@echo "Database available at: localhost:5432"

# Stop all services
down:
	@echo "🛑 Stopping services..."
	docker-compose down
	@echo "✅ Services stopped!"

# Build the Go binary
build:
	@echo "🔨 Building API binary..."
	go build -o bin/finance-api .
	@echo "✅ Build complete! Binary at: bin/finance-api"

# Restart all services
restart:
	@echo "🔄 Restarting services..."
	docker-compose down
	docker-compose up --build -d
	@echo "✅ Services restarted!"

# View logs
logs:
	docker-compose logs -f

# View API logs only
logs-api:
	docker-compose logs -f api

# View database logs only
logs-db:
	docker-compose logs -f postgres

# View web logs only
logs-web:
	docker-compose logs -f web

# Clean up everything
clean:
	@echo "🧹 Cleaning up..."
	docker-compose down -v
	rm -rf api/bin
	@echo "✅ Cleanup complete!"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./... -v
	@echo "✅ Tests complete!"

dev-db:
	@echo "🐘 Starting PostgreSQL database..."
	docker-compose up -d postgres
	@echo "✅ Database started! Accessible at localhost:5432"

# Run API locally (for development)
dev-api:
	@echo "🏃 Running API locally..."
	DB_HOST=localhost go run main.go

# Serve frontend locally with Python
dev-web:
	@echo "🌐 Starting local web server..."
	@echo "Frontend available at: http://localhost:3000"
	cd web && python3 -m http.server 3000

dev: dev-db dev-api dev-web
	@echo "✅ Development environment started!"
	@echo "Press Ctrl+C to stop."
	@wait

# Run database migrations manually
migrate:
	@echo "📊 Running migrations..."
	docker-compose exec api ./finance-api migrate
	@echo "✅ Migrations complete!"

# Open shell in API container
shell-api:
	docker-compose exec api sh

# Open PostgreSQL shell
shell-db:
	docker-compose exec postgres psql -U financeuser -d financedb

# Install Go dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	@echo "✅ Dependencies installed!"

# Format Go code and lintered
fmt:
	@echo "✨ Formatting and Lintering code..."
	go fmt ./...
	go vet ./...
	@echo "✅ Code linted and formatted!"

# Health check
health:
	@echo "🏥 Checking API health..."
	@curl -s http://localhost:8080/api/v1/health || echo "❌ API is not responding"

# Open frontend in browser
open:
	@echo "🌐 Opening frontend in browser..."
	@open http://localhost:3000 2>/dev/null || xdg-open http://localhost:3000 2>/dev/null || echo "Please open http://localhost:3000 in your browser"

# Quick status check
status:
	@echo "📊 Service Status:"
	@docker-compose ps

package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// PostgresDB represents the PostgreSQL database connection
type PostgresDB struct {
	DB *sql.DB
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB() (*PostgresDB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "disable"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &PostgresDB{DB: db}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.DB.Close()
}

// Migrate runs database migrations
func (p *PostgresDB) Migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS assets (
			id UUID PRIMARY KEY,
			type VARCHAR(50) NOT NULL,
			name VARCHAR(255) NOT NULL,
			buy_price DECIMAL(15, 2) NOT NULL,
			current_value DECIMAL(15, 2) NOT NULL,
			currency VARCHAR(10) DEFAULT 'USD',
			quantity DECIMAL(15, 4) NOT NULL,
			purchase_date DATE NOT NULL,
			source VARCHAR(50) DEFAULT 'manual',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS asset_history (
			id UUID PRIMARY KEY,
			asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
			value DECIMAL(15, 2) NOT NULL,
			date DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(asset_id, date)
		)`,
		`CREATE TABLE IF NOT EXISTS debts (
			id UUID PRIMARY KEY,
			type VARCHAR(50) NOT NULL,
			name VARCHAR(255) NOT NULL,
			principal DECIMAL(15, 2) NOT NULL,
			current_value DECIMAL(15, 2) NOT NULL,
			currency VARCHAR(10) DEFAULT 'USD',
			interest_rate DECIMAL(5, 2) DEFAULT 0.00,
			start_date DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(type)`,
		`CREATE INDEX IF NOT EXISTS idx_asset_history_asset_id ON asset_history(asset_id)`,
		`CREATE INDEX IF NOT EXISTS idx_asset_history_date ON asset_history(date)`,
		`CREATE INDEX IF NOT EXISTS idx_debts_type ON debts(type)`,
	}

	for _, migration := range migrations {
		if _, err := p.DB.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}

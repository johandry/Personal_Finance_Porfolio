package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"personal-finance/api/v1/db"
	"personal-finance/api/v1/models"

	"github.com/google/uuid"
)

// parseDate attempts to parse a date string in multiple formats
func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}

	// List of common date formats to try
	formats := []string{
		"2006-01-02",      // ISO format (YYYY-MM-DD)
		"1/2/06",          // M/D/YY
		"01/02/06",        // MM/DD/YY
		"1/2/2006",        // M/D/YYYY
		"01/02/2006",      // MM/DD/YYYY
		"2006/01/02",      // YYYY/MM/DD
		time.RFC3339,      // RFC3339 format
		"Jan 2, 2006",     // Month Day, Year
		"January 2, 2006", // Full month name
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// ExportHandler handles export/import operations
type ExportHandler struct {
	db *db.PostgresDB
}

// NewExportHandler creates a new export handler
func NewExportHandler(database *db.PostgresDB) *ExportHandler {
	return &ExportHandler{db: database}
}

// ExportAssetsJSON handles GET /api/v1/export/assets/json
func (h *ExportHandler) ExportAssetsJSON(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, type, name, buy_price, current_value, currency, quantity, purchase_date, source, created_at, updated_at
		FROM assets
		ORDER BY created_at DESC
	`

	rows, err := h.db.DB.Query(query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch assets")
		return
	}
	defer rows.Close()

	assets := []models.Asset{}
	for rows.Next() {
		var asset models.Asset
		err := rows.Scan(
			&asset.ID, &asset.Type, &asset.Name, &asset.BuyPrice, &asset.CurrentValue,
			&asset.Currency, &asset.Quantity, &asset.PurchaseDate, &asset.Source,
			&asset.CreatedAt, &asset.UpdatedAt,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to parse assets")
			return
		}
		assets = append(assets, asset)
	}

	// Set headers for file download
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=assets_%s.json", time.Now().Format("2006-01-02")))

	json.NewEncoder(w).Encode(assets)
}

// ExportAssetsCSV handles GET /api/v1/export/assets/csv
func (h *ExportHandler) ExportAssetsCSV(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, type, name, buy_price, current_value, currency, quantity, purchase_date, source, created_at, updated_at
		FROM assets
		ORDER BY created_at DESC
	`

	rows, err := h.db.DB.Query(query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch assets")
		return
	}
	defer rows.Close()

	// Set headers for CSV download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=assets_%s.csv", time.Now().Format("2006-01-02")))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write CSV header
	header := []string{"ID", "Type", "Name", "Buy Price", "Current Value", "Currency", "Quantity", "Purchase Date", "Source", "Created At", "Updated At"}
	writer.Write(header)

	// Write data rows
	for rows.Next() {
		var asset models.Asset
		err := rows.Scan(
			&asset.ID, &asset.Type, &asset.Name, &asset.BuyPrice, &asset.CurrentValue,
			&asset.Currency, &asset.Quantity, &asset.PurchaseDate, &asset.Source,
			&asset.CreatedAt, &asset.UpdatedAt,
		)
		if err != nil {
			continue
		}

		row := []string{
			asset.ID,
			string(asset.Type),
			asset.Name,
			fmt.Sprintf("%.2f", asset.BuyPrice),
			fmt.Sprintf("%.2f", asset.CurrentValue),
			asset.Currency,
			fmt.Sprintf("%.4f", asset.Quantity),
			asset.PurchaseDate.Format("2006-01-02"),
			string(asset.Source),
			asset.CreatedAt.Format(time.RFC3339),
			asset.UpdatedAt.Format(time.RFC3339),
		}
		writer.Write(row)
	}
}

// ExportDebtsJSON handles GET /api/v1/export/debts/json
func (h *ExportHandler) ExportDebtsJSON(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, type, name, principal, current_value, currency, interest_rate, start_date, created_at, updated_at
		FROM debts
		ORDER BY created_at DESC
	`

	rows, err := h.db.DB.Query(query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch debts")
		return
	}
	defer rows.Close()

	debts := []models.Debt{}
	for rows.Next() {
		var debt models.Debt
		err := rows.Scan(
			&debt.ID, &debt.Type, &debt.Name, &debt.Principal, &debt.CurrentValue,
			&debt.Currency, &debt.InterestRate, &debt.StartDate,
			&debt.CreatedAt, &debt.UpdatedAt,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to parse debts")
			return
		}
		debts = append(debts, debt)
	}

	// Set headers for file download
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=debts_%s.json", time.Now().Format("2006-01-02")))

	json.NewEncoder(w).Encode(debts)
}

// ExportDebtsCSV handles GET /api/v1/export/debts/csv
func (h *ExportHandler) ExportDebtsCSV(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, type, name, principal, current_value, currency, interest_rate, start_date, created_at, updated_at
		FROM debts
		ORDER BY created_at DESC
	`

	rows, err := h.db.DB.Query(query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch debts")
		return
	}
	defer rows.Close()

	// Set headers for CSV download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=debts_%s.csv", time.Now().Format("2006-01-02")))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write CSV header
	header := []string{"ID", "Type", "Name", "Principal", "Current Value", "Currency", "Interest Rate", "Start Date", "Created At", "Updated At"}
	writer.Write(header)

	// Write data rows
	for rows.Next() {
		var debt models.Debt
		err := rows.Scan(
			&debt.ID, &debt.Type, &debt.Name, &debt.Principal, &debt.CurrentValue,
			&debt.Currency, &debt.InterestRate, &debt.StartDate,
			&debt.CreatedAt, &debt.UpdatedAt,
		)
		if err != nil {
			continue
		}

		row := []string{
			debt.ID,
			string(debt.Type),
			debt.Name,
			fmt.Sprintf("%.2f", debt.Principal),
			fmt.Sprintf("%.2f", debt.CurrentValue),
			debt.Currency,
			fmt.Sprintf("%.2f", debt.InterestRate),
			debt.StartDate.Format("2006-01-02"),
			debt.CreatedAt.Format(time.RFC3339),
			debt.UpdatedAt.Format(time.RFC3339),
		}
		writer.Write(row)
	}
}

// ImportAssetsJSON handles POST /api/v1/import/assets/json
func (h *ExportHandler) ImportAssetsJSON(w http.ResponseWriter, r *http.Request) {
	var assets []models.Asset
	if err := json.NewDecoder(r.Body).Decode(&assets); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	imported := 0
	skipped := 0
	errors := []string{}

	for _, asset := range assets {
		// Generate ID if empty
		if asset.ID == "" {
			asset.ID = uuid.New().String()
		}

		// Check if asset already exists
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM assets WHERE id = $1)`
		h.db.DB.QueryRow(checkQuery, asset.ID).Scan(&exists)

		if exists {
			skipped++
			continue
		}

		// Set timestamps if empty
		if asset.CreatedAt.IsZero() {
			asset.CreatedAt = time.Now()
		}
		if asset.UpdatedAt.IsZero() {
			asset.UpdatedAt = time.Now()
		}

		// Import asset
		query := `
			INSERT INTO assets (id, type, name, buy_price, current_value, currency, quantity, purchase_date, source, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`

		_, err := h.db.DB.Exec(query,
			asset.ID, asset.Type, asset.Name, asset.BuyPrice, asset.CurrentValue,
			asset.Currency, asset.Quantity, asset.PurchaseDate, asset.Source,
			asset.CreatedAt, asset.UpdatedAt,
		)

		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to import %s: %v", asset.Name, err))
			continue
		}

		imported++
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"imported": imported,
		"skipped":  skipped,
		"errors":   errors,
		"total":    len(assets),
	})
}

// ImportAssetsCSV handles POST /api/v1/import/assets/csv
func (h *ExportHandler) ImportAssetsCSV(w http.ResponseWriter, r *http.Request) {
	reader := csv.NewReader(r.Body)
	records, err := reader.ReadAll()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid CSV format")
		return
	}

	if len(records) < 2 {
		respondWithError(w, http.StatusBadRequest, "CSV file is empty")
		return
	}

	imported := 0
	skipped := 0
	errors := []string{}

	// Skip header row
	for i, record := range records[1:] {
		// Minimum required columns: ID, Type, Name, Buy Price, Current Value, Currency, Quantity, Purchase Date, Source
		if len(record) < 9 {
			errors = append(errors, fmt.Sprintf("Row %d: insufficient columns (need at least 9)", i+2))
			continue
		}

		// Generate ID if empty
		assetID := record[0]
		if assetID == "" {
			assetID = uuid.New().String()
		}

		// Check if asset already exists
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM assets WHERE id = $1)`
		h.db.DB.QueryRow(checkQuery, assetID).Scan(&exists)

		if exists {
			skipped++
			continue
		}

		// Parse values
		buyPrice, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: invalid buy price '%s'", i+2, record[3]))
			continue
		}

		currentValue, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: invalid current value '%s'", i+2, record[4]))
			continue
		}

		quantity, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: invalid quantity '%s'", i+2, record[6]))
			continue
		}

		purchaseDate, err := parseDate(record[7])
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: invalid purchase date '%s'", i+2, record[7]))
			continue
		}

		// Parse timestamps if present, otherwise use current time
		var createdAt, updatedAt time.Time
		if len(record) > 9 && record[9] != "" {
			createdAt, _ = time.Parse(time.RFC3339, record[9])
		}
		if len(record) > 10 && record[10] != "" {
			updatedAt, _ = time.Parse(time.RFC3339, record[10])
		}

		// Set timestamps if empty
		if createdAt.IsZero() {
			createdAt = time.Now()
		}
		if updatedAt.IsZero() {
			updatedAt = time.Now()
		}

		// Import asset
		query := `
			INSERT INTO assets (id, type, name, buy_price, current_value, currency, quantity, purchase_date, source, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`

		_, err = h.db.DB.Exec(query,
			assetID, record[1], record[2], buyPrice, currentValue,
			record[5], quantity, purchaseDate, record[8],
			createdAt, updatedAt,
		)

		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: %v", i+2, err))
			continue
		}

		imported++
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"imported": imported,
		"skipped":  skipped,
		"errors":   errors,
		"total":    len(records) - 1,
	})
}

// ImportDebtsJSON handles POST /api/v1/import/debts/json
func (h *ExportHandler) ImportDebtsJSON(w http.ResponseWriter, r *http.Request) {
	var debts []models.Debt
	if err := json.NewDecoder(r.Body).Decode(&debts); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	imported := 0
	skipped := 0
	errors := []string{}

	for _, debt := range debts {
		// Generate ID if empty
		if debt.ID == "" {
			debt.ID = uuid.New().String()
		}

		// Check if debt already exists
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM debts WHERE id = $1)`
		h.db.DB.QueryRow(checkQuery, debt.ID).Scan(&exists)

		if exists {
			skipped++
			continue
		}

		// Set timestamps if empty
		if debt.CreatedAt.IsZero() {
			debt.CreatedAt = time.Now()
		}
		if debt.UpdatedAt.IsZero() {
			debt.UpdatedAt = time.Now()
		}

		// Import debt
		query := `
			INSERT INTO debts (id, type, name, principal, current_value, currency, interest_rate, start_date, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`

		_, err := h.db.DB.Exec(query,
			debt.ID, debt.Type, debt.Name, debt.Principal, debt.CurrentValue,
			debt.Currency, debt.InterestRate, debt.StartDate,
			debt.CreatedAt, debt.UpdatedAt,
		)

		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to import %s: %v", debt.Name, err))
			continue
		}

		imported++
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"imported": imported,
		"skipped":  skipped,
		"errors":   errors,
		"total":    len(debts),
	})
}

// ImportDebtsCSV handles POST /api/v1/import/debts/csv
func (h *ExportHandler) ImportDebtsCSV(w http.ResponseWriter, r *http.Request) {
	reader := csv.NewReader(r.Body)
	records, err := reader.ReadAll()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid CSV format")
		return
	}

	if len(records) < 2 {
		respondWithError(w, http.StatusBadRequest, "CSV file is empty")
		return
	}

	imported := 0
	skipped := 0
	errors := []string{}

	// Skip header row
	for i, record := range records[1:] {
		// Minimum required columns: ID, Type, Name, Principal, Current Value, Currency, Interest Rate, Start Date
		if len(record) < 8 {
			errors = append(errors, fmt.Sprintf("Row %d: insufficient columns (need at least 8)", i+2))
			continue
		}

		// Generate ID if empty
		debtID := record[0]
		if debtID == "" {
			debtID = uuid.New().String()
		}

		// Check if debt already exists
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM debts WHERE id = $1)`
		h.db.DB.QueryRow(checkQuery, debtID).Scan(&exists)

		if exists {
			skipped++
			continue
		}

		// Parse values
		principal, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: invalid principal '%s'", i+2, record[3]))
			continue
		}

		currentValue, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: invalid current value '%s'", i+2, record[4]))
			continue
		}

		interestRate, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: invalid interest rate '%s'", i+2, record[6]))
			continue
		}

		startDate, err := parseDate(record[7])
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: invalid start date '%s'", i+2, record[7]))
			continue
		}

		// Parse timestamps if present, otherwise use current time
		var createdAt, updatedAt time.Time
		if len(record) > 8 && record[8] != "" {
			createdAt, _ = time.Parse(time.RFC3339, record[8])
		}
		if len(record) > 9 && record[9] != "" {
			updatedAt, _ = time.Parse(time.RFC3339, record[9])
		}

		// Set timestamps if empty
		if createdAt.IsZero() {
			createdAt = time.Now()
		}
		if updatedAt.IsZero() {
			updatedAt = time.Now()
		}

		// Import debt
		query := `
			INSERT INTO debts (id, type, name, principal, current_value, currency, interest_rate, start_date, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`

		_, err = h.db.DB.Exec(query,
			debtID, record[1], record[2], principal, currentValue,
			record[5], interestRate, startDate,
			createdAt, updatedAt,
		)

		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: %v", i+2, err))
			continue
		}

		imported++
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"imported": imported,
		"skipped":  skipped,
		"errors":   errors,
		"total":    len(records) - 1,
	})
}

// ExportAll handles GET /api/v1/export/all/json
func (h *ExportHandler) ExportAll(w http.ResponseWriter, r *http.Request) {
	// Fetch all assets
	assetsQuery := `SELECT id, type, name, buy_price, current_value, currency, quantity, purchase_date, source, created_at, updated_at FROM assets ORDER BY created_at DESC`
	assetsRows, err := h.db.DB.Query(assetsQuery)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch assets")
		return
	}
	defer assetsRows.Close()

	assets := []models.Asset{}
	for assetsRows.Next() {
		var asset models.Asset
		assetsRows.Scan(&asset.ID, &asset.Type, &asset.Name, &asset.BuyPrice, &asset.CurrentValue,
			&asset.Currency, &asset.Quantity, &asset.PurchaseDate, &asset.Source,
			&asset.CreatedAt, &asset.UpdatedAt)
		assets = append(assets, asset)
	}

	// Fetch all debts
	debtsQuery := `SELECT id, type, name, principal, current_value, currency, interest_rate, start_date, created_at, updated_at FROM debts ORDER BY created_at DESC`
	debtsRows, err := h.db.DB.Query(debtsQuery)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch debts")
		return
	}
	defer debtsRows.Close()

	debts := []models.Debt{}
	for debtsRows.Next() {
		var debt models.Debt
		debtsRows.Scan(&debt.ID, &debt.Type, &debt.Name, &debt.Principal, &debt.CurrentValue,
			&debt.Currency, &debt.InterestRate, &debt.StartDate,
			&debt.CreatedAt, &debt.UpdatedAt)
		debts = append(debts, debt)
	}

	exportData := map[string]interface{}{
		"assets":      assets,
		"debts":       debts,
		"exported_at": time.Now(),
		"version":     "1.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=portfolio_%s.json", time.Now().Format("2006-01-02")))

	json.NewEncoder(w).Encode(exportData)
}

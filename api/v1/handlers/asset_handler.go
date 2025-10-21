package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"personal-finance/api/v1/db"
	"personal-finance/api/v1/models"
	"personal-finance/api/v1/services"
)

// AssetHandler handles asset-related requests
type AssetHandler struct {
	db         *db.PostgresDB
	marketData *services.MarketDataService
}

// NewAssetHandler creates a new asset handler
func NewAssetHandler(database *db.PostgresDB, marketDataService *services.MarketDataService) *AssetHandler {
	return &AssetHandler{
		db:         database,
		marketData: marketDataService,
	}
}

// CreateAsset handles POST /api/v1/assets
func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate required fields
	if req.Name == "" || req.Type == "" || req.BuyPrice <= 0 || req.Quantity <= 0 {
		respondWithError(w, http.StatusBadRequest, "Missing or invalid required fields")
		return
	}

	purchaseDate, err := time.Parse("2006-01-02", req.PurchaseDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid purchase_date format (use YYYY-MM-DD)")
		return
	}

	// Set current value to buy price if not provided
	currentValue := req.BuyPrice
	if req.CurrentValue != nil {
		currentValue = *req.CurrentValue
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}

	asset := models.Asset{
		ID:           uuid.New().String(),
		Type:         req.Type,
		Name:         req.Name,
		BuyPrice:     req.BuyPrice,
		CurrentValue: currentValue,
		Currency:     req.Currency,
		Quantity:     req.Quantity,
		PurchaseDate: purchaseDate,
		Source:       req.Source,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query := `
		INSERT INTO assets (id, type, name, buy_price, current_value, currency, quantity, purchase_date, source, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = h.db.DB.Exec(query,
		asset.ID, asset.Type, asset.Name, asset.BuyPrice, asset.CurrentValue,
		asset.Currency, asset.Quantity, asset.PurchaseDate, asset.Source,
		asset.CreatedAt, asset.UpdatedAt,
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create asset")
		return
	}

	// Create initial history entry
	h.addHistoryEntry(asset.ID, asset.CurrentValue, time.Now())

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"id":            asset.ID,
		"name":          asset.Name,
		"current_value": asset.CurrentValue,
		"profit_loss":   asset.ProfitLoss(),
	})
}

// ListAssets handles GET /api/v1/assets
func (h *AssetHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
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

		// Fetch real-time price for stocks
		if string(asset.Type) == "stock" && string(asset.Source) == "market_api" {
			currentPrice, err := h.marketData.GetCurrentValue(
				string(asset.Type),
				asset.Name,
				asset.CurrentValue,
				string(asset.Source),
			)
			if err == nil {
				asset.CurrentValue = currentPrice
			}
			// If error, keep the stored value
		}

		assets = append(assets, asset)
	}

	respondWithJSON(w, http.StatusOK, assets)
}

// GetAsset handles GET /api/v1/assets/{id}
func (h *AssetHandler) GetAsset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := `
		SELECT id, type, name, buy_price, current_value, currency, quantity, purchase_date, source, created_at, updated_at
		FROM assets
		WHERE id = $1
	`

	var asset models.Asset
	err := h.db.DB.QueryRow(query, id).Scan(
		&asset.ID, &asset.Type, &asset.Name, &asset.BuyPrice, &asset.CurrentValue,
		&asset.Currency, &asset.Quantity, &asset.PurchaseDate, &asset.Source,
		&asset.CreatedAt, &asset.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Asset not found")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch asset")
		return
	}

	// Fetch real-time price for stocks
	if string(asset.Type) == "stock" && string(asset.Source) == "market_api" {
		currentPrice, err := h.marketData.GetCurrentValue(
			string(asset.Type),
			asset.Name,
			asset.CurrentValue,
			string(asset.Source),
		)
		if err == nil {
			asset.CurrentValue = currentPrice
		}
		// If error, keep the stored value
	}

	respondWithJSON(w, http.StatusOK, asset)
}

// UpdateAsset handles PUT /api/v1/assets/{id}
func (h *AssetHandler) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Build dynamic update query
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.CurrentValue != nil {
		updates["current_value"] = *req.CurrentValue
	}
	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.Source != nil {
		updates["source"] = *req.Source
	}

	if len(updates) == 0 {
		respondWithError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	updates["updated_at"] = time.Now()

	// Execute update
	query := "UPDATE assets SET "
	args := []interface{}{}
	i := 1
	for key, val := range updates {
		if i > 1 {
			query += ", "
		}
		query += key + " = $" + string(rune(i+48))
		args = append(args, val)
		i++
	}
	query += " WHERE id = $" + string(rune(i+48))
	args = append(args, id)

	result, err := h.db.DB.Exec(query, args...)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update asset")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Asset not found")
		return
	}

	// Add history entry if current value was updated
	if req.CurrentValue != nil {
		h.addHistoryEntry(id, *req.CurrentValue, time.Now())
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Asset updated successfully"})
}

// DeleteAsset handles DELETE /api/v1/assets/{id}
func (h *AssetHandler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := "DELETE FROM assets WHERE id = $1"
	result, err := h.db.DB.Exec(query, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete asset")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Asset not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Asset deleted successfully"})
}

// GetAssetHistory handles GET /api/v1/assets/{id}/history
func (h *AssetHandler) GetAssetHistory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := `
		SELECT id, asset_id, value, date, created_at
		FROM asset_history
		WHERE asset_id = $1
		ORDER BY date DESC
	`

	rows, err := h.db.DB.Query(query, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch asset history")
		return
	}
	defer rows.Close()

	history := []models.AssetHistory{}
	for rows.Next() {
		var h models.AssetHistory
		err := rows.Scan(&h.ID, &h.AssetID, &h.Value, &h.Date, &h.CreatedAt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to parse history")
			return
		}
		history = append(history, h)
	}

	respondWithJSON(w, http.StatusOK, history)
}

// addHistoryEntry adds a history entry for an asset
func (h *AssetHandler) addHistoryEntry(assetID string, value float64, date time.Time) error {
	query := `
		INSERT INTO asset_history (id, asset_id, value, date, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (asset_id, date) DO UPDATE SET value = $3
	`

	_, err := h.db.DB.Exec(query, uuid.New().String(), assetID, value, date.Format("2006-01-02"), time.Now())
	return err
}

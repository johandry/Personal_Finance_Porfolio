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
)

// DebtHandler handles debt-related requests
type DebtHandler struct {
	db *db.PostgresDB
}

// NewDebtHandler creates a new debt handler
func NewDebtHandler(database *db.PostgresDB) *DebtHandler {
	return &DebtHandler{db: database}
}

// CreateDebt handles POST /api/v1/debts
func (h *DebtHandler) CreateDebt(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDebtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate required fields
	if req.Name == "" || req.Type == "" || req.Principal <= 0 {
		respondWithError(w, http.StatusBadRequest, "Missing or invalid required fields")
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid start_date format (use YYYY-MM-DD)")
		return
	}

	// Set current value to principal if not provided
	currentValue := req.Principal
	if req.CurrentValue != nil {
		currentValue = *req.CurrentValue
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}

	debt := models.Debt{
		ID:           uuid.New().String(),
		Type:         req.Type,
		Name:         req.Name,
		Principal:    req.Principal,
		CurrentValue: currentValue,
		Currency:     req.Currency,
		InterestRate: req.InterestRate,
		StartDate:    startDate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query := `
		INSERT INTO debts (id, type, name, principal, current_value, currency, interest_rate, start_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = h.db.DB.Exec(query,
		debt.ID, debt.Type, debt.Name, debt.Principal, debt.CurrentValue,
		debt.Currency, debt.InterestRate, debt.StartDate,
		debt.CreatedAt, debt.UpdatedAt,
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create debt")
		return
	}

	respondWithJSON(w, http.StatusCreated, debt)
}

// ListDebts handles GET /api/v1/debts
func (h *DebtHandler) ListDebts(w http.ResponseWriter, r *http.Request) {
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

	respondWithJSON(w, http.StatusOK, debts)
}

// GetDebt handles GET /api/v1/debts/{id}
func (h *DebtHandler) GetDebt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := `
		SELECT id, type, name, principal, current_value, currency, interest_rate, start_date, created_at, updated_at
		FROM debts
		WHERE id = $1
	`

	var debt models.Debt
	err := h.db.DB.QueryRow(query, id).Scan(
		&debt.ID, &debt.Type, &debt.Name, &debt.Principal, &debt.CurrentValue,
		&debt.Currency, &debt.InterestRate, &debt.StartDate,
		&debt.CreatedAt, &debt.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Debt not found")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch debt")
		return
	}

	respondWithJSON(w, http.StatusOK, debt)
}

// UpdateDebt handles PUT /api/v1/debts/{id}
func (h *DebtHandler) UpdateDebt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateDebtRequest
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
	if req.InterestRate != nil {
		updates["interest_rate"] = *req.InterestRate
	}

	if len(updates) == 0 {
		respondWithError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	updates["updated_at"] = time.Now()

	// Execute update
	query := "UPDATE debts SET "
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
		respondWithError(w, http.StatusInternalServerError, "Failed to update debt")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Debt not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Debt updated successfully"})
}

// DeleteDebt handles DELETE /api/v1/debts/{id}
func (h *DebtHandler) DeleteDebt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := "DELETE FROM debts WHERE id = $1"
	result, err := h.db.DB.Exec(query, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete debt")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Debt not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Debt deleted successfully"})
}

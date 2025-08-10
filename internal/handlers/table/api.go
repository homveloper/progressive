package table

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// APIHandler handles table data API related requests
type APIHandler struct {
	db *sqlx.DB
}

// NewAPIHandler creates a new APIHandler instance
func NewAPIHandler(db *sqlx.DB) *APIHandler {
	return &APIHandler{db: db}
}

// DataHandler handles table data API requests with pagination
func (h *APIHandler) DataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract table ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/table/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Table ID required", http.StatusBadRequest)
		return
	}
	tableID := parts[0]

	// Parse pagination parameters
	page := parseInt(r.URL.Query().Get("page"), 1)
	limit := parseInt(r.URL.Query().Get("limit"), 20)
	offset := (page - 1) * limit

	// Get table metadata
	query := `
		SELECT id, name, description, schema, record_count, created_at, updated_at 
		FROM tables 
		WHERE id = $1
	`

	var table struct {
		ID          string          `db:"id"`
		Name        string          `db:"name"`
		Description string          `db:"description"`
		Schema      json.RawMessage `db:"schema"`
		RecordCount int             `db:"record_count"`
		CreatedAt   time.Time       `db:"created_at"`
		UpdatedAt   time.Time       `db:"updated_at"`
	}

	err := h.db.Get(&table, query, tableID)
	if err != nil {
		http.Error(w, "Table not found", http.StatusNotFound)
		return
	}

	// Get paginated records for this table
	recordsQuery := `
		SELECT id, data, created_at 
		FROM records 
		WHERE table_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	rows, err := h.db.Query(recordsQuery, tableID, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch records: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var records []map[string]interface{}
	for rows.Next() {
		var id int
		var data json.RawMessage
		var createdAt time.Time

		if err := rows.Scan(&id, &data, &createdAt); err == nil {
			record := make(map[string]interface{})
			if err := json.Unmarshal(data, &record); err == nil {
				record["_id"] = id
				record["_created_at"] = createdAt
				records = append(records, record)
			}
		}
	}

	// Create response
	response := map[string]interface{}{
		"table": map[string]interface{}{
			"id":           table.ID,
			"name":         table.Name,
			"description":  table.Description,
			"schema":       json.RawMessage(table.Schema),
			"record_count": table.RecordCount,
			"created_at":   table.CreatedAt,
			"updated_at":   table.UpdatedAt,
		},
		"records": records,
		"pagination": map[string]interface{}{
			"page":     page,
			"limit":    limit,
			"total":    table.RecordCount,
			"has_more": (offset + len(records)) < table.RecordCount,
		},
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON response: %v", err), http.StatusInternalServerError)
		return
	}
}

// RecordHandler handles individual record operations
func (h *APIHandler) RecordHandler(w http.ResponseWriter, r *http.Request) {
	// Extract table ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/table/")
	parts := strings.Split(path, "/")
	if len(parts) < 3 || parts[0] == "" || parts[1] != "record" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	tableID := parts[0]

	switch r.Method {
	case "POST":
		h.createRecord(w, r, tableID)
	case "PATCH":
		if len(parts) < 3 || parts[2] == "" {
			http.Error(w, "Record ID required for PATCH", http.StatusBadRequest)
			return
		}
		recordID := parts[2]
		h.updateRecord(w, r, tableID, recordID)
	case "DELETE":
		if len(parts) < 3 || parts[2] == "" {
			http.Error(w, "Record ID required for DELETE", http.StatusBadRequest)
			return
		}
		recordID := parts[2]
		h.deleteRecord(w, r, tableID, recordID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ImportHandler handles data import requests
func (h *APIHandler) ImportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract table ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/table/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] != "import" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	tableID := parts[0]

	// Parse request body
	var importRequest struct {
		Data []map[string]interface{} `json:"data"`
		Mode string                   `json:"mode"` // "replace" or "append"
	}

	if err := json.NewDecoder(r.Body).Decode(&importRequest); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	if len(importRequest.Data) == 0 {
		http.Error(w, "No data to import", http.StatusBadRequest)
		return
	}

	// Validate mode
	if importRequest.Mode != "replace" && importRequest.Mode != "append" {
		importRequest.Mode = "replace" // default
	}

	// Start transaction
	tx, err := h.db.Beginx()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start transaction: %v", err), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// If replace mode, delete existing records
	if importRequest.Mode == "replace" {
		_, err := tx.Exec("DELETE FROM records WHERE table_id = $1", tableID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to clear existing records: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Insert new records
	insertQuery := `INSERT INTO records (table_id, data, created_at) VALUES ($1, $2, $3)`
	importedCount := 0
	now := time.Now()

	for _, record := range importRequest.Data {
		// Convert record to JSON
		recordJSON, err := json.Marshal(record)
		if err != nil {
			continue // Skip invalid records
		}

		// Insert record
		_, err = tx.Exec(insertQuery, tableID, json.RawMessage(recordJSON), now)
		if err == nil {
			importedCount++
		}
	}

	if importedCount == 0 {
		http.Error(w, "No valid records could be imported", http.StatusBadRequest)
		return
	}

	// Update table record count
	countQuery := `UPDATE tables SET record_count = (SELECT COUNT(*) FROM records WHERE table_id = $1), updated_at = $2 WHERE id = $1`
	_, err = tx.Exec(countQuery, tableID, now)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update table count: %v", err), http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to commit transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"imported": importedCount,
		"mode":     importRequest.Mode,
	})
}

// ExportHandler handles data export requests
func (h *APIHandler) ExportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract table ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/table/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] != "export" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	tableID := parts[0]

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	// Get all records for the table
	recordsQuery := `
		SELECT data FROM records 
		WHERE table_id = $1 
		ORDER BY created_at DESC
	`

	rows, err := h.db.Query(recordsQuery, tableID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch records: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var records []map[string]interface{}
	for rows.Next() {
		var data json.RawMessage
		if err := rows.Scan(&data); err == nil {
			record := make(map[string]interface{})
			if err := json.Unmarshal(data, &record); err == nil {
				records = append(records, record)
			}
		}
	}

	// Export based on format
	switch format {
	case "json":
		h.exportJSON(w, records)
	case "csv":
		h.exportCSV(w, records)
	case "excel":
		h.exportExcel(w, records)
	default:
		http.Error(w, "Unsupported format", http.StatusBadRequest)
	}
}

// createRecord creates a new record
func (h *APIHandler) createRecord(w http.ResponseWriter, r *http.Request, tableID string) {
	var record map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Insert record into database
	data, _ := json.Marshal(record)
	query := `INSERT INTO records (table_id, data, created_at) VALUES ($1, $2, $3) RETURNING id`

	var recordID int
	err := h.db.QueryRow(query, tableID, json.RawMessage(data), time.Now()).Scan(&recordID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create record: %v", err), http.StatusInternalServerError)
		return
	}

	// Update table record count
	h.updateRecordCount(tableID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      recordID,
	})
}

// updateRecord updates an existing record
func (h *APIHandler) updateRecord(w http.ResponseWriter, r *http.Request, tableID, recordID string) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get existing record
	var existingData json.RawMessage
	query := `SELECT data FROM records WHERE id = $1 AND table_id = $2`
	err := h.db.QueryRow(query, recordID, tableID).Scan(&existingData)
	if err != nil {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	// Merge updates with existing data
	var record map[string]interface{}
	json.Unmarshal(existingData, &record)

	for key, value := range updates {
		record[key] = value
	}

	// Update record in database
	newData, _ := json.Marshal(record)
	updateQuery := `UPDATE records SET data = $1, updated_at = $2 WHERE id = $3 AND table_id = $4`
	_, err = h.db.Exec(updateQuery, json.RawMessage(newData), time.Now(), recordID, tableID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update record: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// deleteRecord deletes a record
func (h *APIHandler) deleteRecord(w http.ResponseWriter, r *http.Request, tableID, recordID string) {
	query := `DELETE FROM records WHERE id = $1 AND table_id = $2`
	result, err := h.db.Exec(query, recordID, tableID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete record: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	// Update table record count
	h.updateRecordCount(tableID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// updateRecordCount updates the record count for a table
func (h *APIHandler) updateRecordCount(tableID string) {
	query := `UPDATE tables SET record_count = (SELECT COUNT(*) FROM records WHERE table_id = $1) WHERE id = $1`
	h.db.Exec(query, tableID)
}

// Export functions
func (h *APIHandler) exportJSON(w http.ResponseWriter, records []map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=export.json")
	json.NewEncoder(w).Encode(records)
}

func (h *APIHandler) exportCSV(w http.ResponseWriter, records []map[string]interface{}) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=export.csv")

	if len(records) == 0 {
		return
	}

	// Get headers from first record
	var headers []string
	for key := range records[0] {
		headers = append(headers, key)
	}

	// Write CSV header
	for i, header := range headers {
		if i > 0 {
			w.Write([]byte(","))
		}
		w.Write([]byte(fmt.Sprintf("\"%s\"", header)))
	}
	w.Write([]byte("\n"))

	// Write CSV rows
	for _, record := range records {
		for i, header := range headers {
			if i > 0 {
				w.Write([]byte(","))
			}
			value := fmt.Sprintf("%v", record[header])
			w.Write([]byte(fmt.Sprintf("\"%s\"", strings.ReplaceAll(value, "\"", "\"\""))))
		}
		w.Write([]byte("\n"))
	}
}

func (h *APIHandler) exportExcel(w http.ResponseWriter, records []map[string]interface{}) {
	// For now, just export as CSV with Excel content type
	w.Header().Set("Content-Type", "application/vnd.ms-excel")
	w.Header().Set("Content-Disposition", "attachment; filename=export.xlsx")
	h.exportCSV(w, records)
}

// parseInt parses string to int with default value
func parseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	if val := parseIntHelper(s); val > 0 {
		return val
	}
	return defaultValue
}

func parseIntHelper(s string) int {
	result := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result = result*10 + int(r-'0')
		} else {
			return 0
		}
	}
	return result
}
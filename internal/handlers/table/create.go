package table

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"progressive/internal/models"
	"progressive/internal/pages"
	"time"

	"github.com/jmoiron/sqlx"
)

// CreateHandler handles table creation related requests
type CreateHandler struct {
	db *sqlx.DB
}

// NewCreateHandler creates a new CreateHandler instance
func NewCreateHandler(db *sqlx.DB) *CreateHandler {
	return &CreateHandler{db: db}
}

// PageHandler renders the table creation page (GET only)
func (h *CreateHandler) PageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get templates data from models
	templates := models.GetTableTemplates()

	// Render the table creation page
	err := pages.TableCreatePage(templates).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// APIHandler handles table creation API requests (POST only)
func (h *CreateHandler) APIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.handleTableCreationAPI(w, r)
}

// TableCreateRequest represents the JSON request for table creation
type TableCreateRequest struct {
	TableName  string `json:"table_name"`
	Schema     string `json:"schema"`
	DataOption string `json:"data_option"`
}

// handleTableCreationAPI processes the table creation API request
func (h *CreateHandler) handleTableCreationAPI(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var req TableCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON request: %v", err), http.StatusBadRequest)
		return
	}

	// Extract table data from request
	tableName := req.TableName
	schemaJSON := req.Schema
	dataOption := req.DataOption

	// Debug: Log received JSON values
	log.Printf("🔍 JSON data received - table_name: '%s', schema length: %d, data_option: '%s'",
		tableName, len(schemaJSON), dataOption)

	// Validate required fields
	if tableName == "" {
		http.Error(w, "Table name is required", http.StatusBadRequest)
		return
	}

	if schemaJSON == "" {
		http.Error(w, "Schema is required", http.StatusBadRequest)
		return
	}

	// Validate JSON schema syntax
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(schemaJSON), &schema); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON schema syntax: %v", err), http.StatusBadRequest)
		return
	}

	// Validate schema structure - must be object type
	schemaType, typeExists := schema["type"]
	if !typeExists {
		http.Error(w, "Schema must have a 'type' field", http.StatusBadRequest)
		return
	}

	if schemaType != "object" {
		http.Error(w, fmt.Sprintf("Schema type must be 'object', got '%v'", schemaType), http.StatusBadRequest)
		return
	}

	// Validate properties exist
	propertiesField, propertiesExists := schema["properties"]
	if !propertiesExists {
		http.Error(w, "Schema must contain 'properties' field", http.StatusBadRequest)
		return
	}

	properties, ok := propertiesField.(map[string]interface{})
	if !ok {
		http.Error(w, "Schema 'properties' field must be an object", http.StatusBadRequest)
		return
	}

	if len(properties) == 0 {
		http.Error(w, "Schema must contain at least one property", http.StatusBadRequest)
		return
	}

	// Save table to database
	tableID := generateTableID(tableName)

	// Create table in database
	query := `
		INSERT INTO tables (id, name, description, schema, record_count, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	description := "사용자가 생성한 테이블: " + tableName

	_, err := h.db.Exec(query, tableID, tableName, description, json.RawMessage(schemaJSON), 0, now, now)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save table to database: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Table created successfully: %s (ID: %s)", tableName, tableID)

	// Create response
	response := map[string]interface{}{
		"success":    true,
		"tableId":    tableID,
		"name":       tableName,
		"schema":     schema,
		"dataOption": dataOption,
		"redirect":   "/table/" + tableID,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Send JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON response: %v", err), http.StatusInternalServerError)
		return
	}
}

// generateTableID creates a unique table ID
func generateTableID(name string) string {
	// Create a more unique ID with timestamp and random suffix
	timestamp := time.Now().UnixNano()
	randomSuffix := rand.Intn(10000)
	safeName := sanitizeForID(name)
	if len(safeName) > 20 {
		safeName = safeName[:20]
	}
	return fmt.Sprintf("table_%s_%d_%d", safeName, timestamp, randomSuffix)
}

// sanitizeForID creates a URL-safe ID from a string
func sanitizeForID(s string) string {
	// Simple sanitization - in production, use a more robust method
	result := ""
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result += string(r)
		} else if r == ' ' {
			result += "_"
		}
	}
	return result
}

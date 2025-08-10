package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"progressive/internal/domain/schematemplate/repository"
	"progressive/internal/pages"

	"github.com/jmoiron/sqlx"
)

// Handlers holds all HTTP handlers with database dependency
type Handlers struct {
	db           *sqlx.DB
	templateRepo repository.SchemaTemplateRepository
	Table        *TableHandlers
}

// NewHandlers creates a new Handlers instance with database
func NewHandlers(db *sqlx.DB) *Handlers {
	// Initialize repository with default templates
	templateRepo, err := repository.NewPostgresRepositoryWithDefaults(context.Background(), db)
	if err != nil {
		log.Printf("Failed to initialize template repository with defaults: %v", err)
		// Fallback to basic repository without defaults
		templateRepo = repository.NewPostgresRepository(db)
	} else {
		log.Println("✅ Default templates initialized successfully")
	}

	return &Handlers{
		db:           db,
		templateRepo: templateRepo,
		Table:        NewTableHandlers(db),
	}
}

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
	component := pages.HomePage()
	component.Render(r.Context(), w)
}

func (h *Handlers) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	component := pages.Dashboard()
	component.Render(r.Context(), w)
}

// TemplatesAPIHandler returns all templates as JSON
func (h *Handlers) TemplatesAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	templates, err := h.templateRepo.FindAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

// TablesAPIHandler handles table CRUD operations
func (h *Handlers) TablesAPIHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getTablesHandler(w, r)
	case "POST":
		h.createTableHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) getTablesHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, name, description, schema, record_count, created_at, updated_at 
		FROM tables 
		ORDER BY updated_at DESC
	`

	var tables []map[string]interface{}
	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		table := make(map[string]interface{})
		var schemaJSON json.RawMessage
		var id, name, description string
		var recordCount int
		var createdAt, updatedAt interface{}

		err := rows.Scan(&id, &name, &description, &schemaJSON, &recordCount, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		table["id"] = id
		table["name"] = name
		table["description"] = description
		table["schema"] = schemaJSON
		table["record_count"] = recordCount
		table["created_at"] = createdAt
		table["updated_at"] = updatedAt

		tables = append(tables, table)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tables)
}

func (h *Handlers) createTableHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID          string          `json:"id"`
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Schema      json.RawMessage `json:"schema"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO tables (id, name, description, schema) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id string
	err := h.db.QueryRow(query, payload.ID, payload.Name, payload.Description, payload.Schema).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// TableAPIHandler handles single table operations
func (h *Handlers) TableAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Extract table ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/table/")
	tableID := strings.TrimSuffix(path, "/")

	if tableID == "" {
		http.Error(w, "Table ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		h.getTableHandler(w, r, tableID)
	case "PUT":
		h.updateTableHandler(w, r, tableID)
	case "DELETE":
		h.deleteTableHandler(w, r, tableID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) getTableHandler(w http.ResponseWriter, r *http.Request, tableID string) {
	query := `
		SELECT id, name, description, schema, record_count, created_at, updated_at 
		FROM tables 
		WHERE id = $1
	`

	table := make(map[string]interface{})
	var schemaJSON json.RawMessage
	var id, name, description string
	var recordCount int
	var createdAt, updatedAt interface{}

	err := h.db.QueryRow(query, tableID).Scan(&id, &name, &description, &schemaJSON, &recordCount, &createdAt, &updatedAt)

	if err != nil {
		http.Error(w, "Table not found", http.StatusNotFound)
		return
	}

	table["id"] = id
	table["name"] = name
	table["description"] = description
	table["schema"] = schemaJSON
	table["record_count"] = recordCount
	table["created_at"] = createdAt
	table["updated_at"] = updatedAt

	// Get records for this table
	recordsQuery := `SELECT id, data FROM records WHERE table_id = $1 ORDER BY created_at DESC`
	rows, err := h.db.Query(recordsQuery, tableID)
	if err == nil {
		defer rows.Close()

		var records []map[string]interface{}
		for rows.Next() {
			var id int
			var data json.RawMessage
			if err := rows.Scan(&id, &data); err == nil {
				record := make(map[string]interface{})
				json.Unmarshal(data, &record)
				record["_id"] = id
				records = append(records, record)
			}
		}
		table["records"] = records
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(table)
}

func (h *Handlers) updateTableHandler(w http.ResponseWriter, r *http.Request, tableID string) {
	var payload struct {
		Records []json.RawMessage `json:"records"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := h.db.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete existing records
	if _, err := tx.Exec("DELETE FROM records WHERE table_id = $1", tableID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert new records
	for _, record := range payload.Records {
		_, err := tx.Exec("INSERT INTO records (table_id, data) VALUES ($1, $2)", tableID, record)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handlers) deleteTableHandler(w http.ResponseWriter, r *http.Request, tableID string) {
	_, err := h.db.Exec("DELETE FROM tables WHERE id = $1", tableID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

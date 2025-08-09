package handlers

import (
	"encoding/json"
	"net/http"
	"progressive/internal/pages"
	"progressive/internal/models"
)

// TableCreateHandler renders the table creation page
func TableCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Get templates data
		templates := models.GetTableTemplates()
		
		// Render the table creation page
		err := pages.TableCreatePage(templates).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "POST":
		// Handle table creation
		handleTableCreation(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTableCreation processes the table creation request
func handleTableCreation(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract table data from form
	tableName := r.FormValue("table_name")
	schemaJSON := r.FormValue("schema")
	dataOption := r.FormValue("data_option")

	if tableName == "" || schemaJSON == "" {
		http.Error(w, "Table name and schema are required", http.StatusBadRequest)
		return
	}

	// Validate JSON schema
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(schemaJSON), &schema); err != nil {
		http.Error(w, "Invalid JSON schema: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate schema structure
	if schema["type"] != "object" {
		http.Error(w, "Schema must be of type 'object'", http.StatusBadRequest)
		return
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok || len(properties) == 0 {
		http.Error(w, "Schema must contain properties", http.StatusBadRequest)
		return
	}

	// TODO: Save table to database
	// For now, we'll create a mock table ID and redirect to the table editor

	// Generate a mock table ID (in real implementation, this would come from the database)
	tableID := generateTableID(tableName)

	// Create response
	response := map[string]interface{}{
		"success":  true,
		"tableId":  tableID,
		"name":     tableName,
		"schema":   schema,
		"dataOption": dataOption,
		"redirect": "/table/" + tableID,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Send JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// generateTableID creates a simple table ID (in real implementation, use UUID or database auto-increment)
func generateTableID(name string) string {
	// For demo purposes, create a simple ID based on the name
	// In production, use proper UUID generation
	return "table_" + sanitizeForID(name)
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

// TableEditorHandler renders the table editor page
func TableEditorHandler(w http.ResponseWriter, r *http.Request) {
	// Extract table ID from URL path
	// For now, we'll create a basic editor page
	err := pages.TableEditorPage().Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
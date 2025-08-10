package table

import (
	"net/http"
	"progressive/internal/pages"
	"strings"

	"github.com/jmoiron/sqlx"
)

// EditorHandler handles table editor related requests
type EditorHandler struct {
	db *sqlx.DB
}

// NewEditorHandler creates a new EditorHandler instance
func NewEditorHandler(db *sqlx.DB) *EditorHandler {
	return &EditorHandler{db: db}
}

// PageHandler renders the table editor page (GET only)
func (h *EditorHandler) PageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract table ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/table/")
	tableID := strings.TrimSuffix(path, "/")

	if tableID == "" || tableID == "table" {
		http.Error(w, "Table ID required", http.StatusBadRequest)
		return
	}

	// Render the table editor page (data will be fetched via API)
	err := pages.TableEditorPage().Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
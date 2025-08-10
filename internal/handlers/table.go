package handlers

import (
	"net/http"
	"progressive/internal/handlers/table"

	"github.com/jmoiron/sqlx"
)

// TableHandlers holds all table-related handlers
type TableHandlers struct {
	Create *table.CreateHandler
	Editor *table.EditorHandler
	API    *table.APIHandler
}

// NewTableHandlers creates a new TableHandlers instance
func NewTableHandlers(db *sqlx.DB) *TableHandlers {
	return &TableHandlers{
		Create: table.NewCreateHandler(db),
		Editor: table.NewEditorHandler(db),
		API:    table.NewAPIHandler(db),
	}
}

// Legacy handlers for backward compatibility
func (h *Handlers) TableCreatePageHandler(w http.ResponseWriter, r *http.Request) {
	h.Table.Create.PageHandler(w, r)
}

func (h *Handlers) TableCreateAPIHandler(w http.ResponseWriter, r *http.Request) {
	h.Table.Create.APIHandler(w, r)
}

func (h *Handlers) TableEditorPageHandler(w http.ResponseWriter, r *http.Request) {
	h.Table.Editor.PageHandler(w, r)
}

func (h *Handlers) TableDataAPIHandler(w http.ResponseWriter, r *http.Request) {
	h.Table.API.DataHandler(w, r)
}
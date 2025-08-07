package editor

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"progressive/spreadsheet/components/grid"
	"progressive/spreadsheet/models/data"
	"progressive/spreadsheet/services/storage"
	"progressive/spreadsheet/utils"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// EditorState represents the current state of the editor
type EditorState struct {
	Spreadsheet  *data.Spreadsheet     `json:"spreadsheet"`
	Schema       *utils.JSONSchema     `json:"schema,omitempty"`
	Columns      []grid.ColumnMetadata `json:"columns"`
	Data         map[string]string     `json:"data"`
	SelectedCell string                `json:"selected_cell"`
	IsEditing    bool                  `json:"is_editing"`
	EditValue    string                `json:"edit_value"`
	LastSaved    time.Time             `json:"last_saved"`
	HasChanges   bool                  `json:"has_changes"`
	Error        string                `json:"error"`
}

// Editor represents the main Excel-like spreadsheet editor page.
// It provides comprehensive functionality for JSON schema-based data editing
// with real-time validation, auto-save, and Excel-like user experience.
//
// Key Features:
//   - JSON schema file upload and parsing
//   - DatabaseTableGrid integration for metadata display
//   - Cell editing with type-specific validation
//   - Keyboard navigation (Enter, Tab, Arrow keys)
//   - Save/Export functionality with local storage
//   - Auto-save with configurable intervals
//   - Status bar showing current cell information
//   - Responsive design for multiple screen sizes
//
// Example usage:
//
//	editor := &Editor{
//	    AutoSaveInterval: 30 * time.Second,
//	    MaxRows:          1000,
//	    MaxCols:          50,
//	}
type Editor struct {
	app.Compo

	// Configuration
	AutoSaveInterval time.Duration
	MaxRows          int
	MaxCols          int
	StoragePrefix    string

	// State
	state   EditorState
	storage *storage.LocalStorage

	// UI State
	showFileUpload bool
	autoSaveTimer  app.Value
}

// OnMount initializes the Editor component when mounted to the DOM.
// Sets up default configuration, initializes storage service, and attempts
// to restore the last session state from localStorage.
func (e *Editor) OnMount(ctx app.Context) {
	// Initialize configuration with sensible defaults
	if e.AutoSaveInterval == 0 {
		e.AutoSaveInterval = 30 * time.Second
	}
	if e.MaxRows == 0 {
		e.MaxRows = 1000
	}
	if e.MaxCols == 0 {
		e.MaxCols = 50
	}
	if e.StoragePrefix == "" {
		e.StoragePrefix = "excel_editor"
	}

	// Initialize storage service
	e.storage = storage.NewLocalStorage(e.StoragePrefix)

	// Initialize state with default columns for empty spreadsheet
	e.state = EditorState{
		Data:         make(map[string]string),
		Columns:      e.createDefaultColumns(),
		SelectedCell: "0-0",
		LastSaved:    time.Now(),
	}

	// Try to restore last session
	e.restoreLastSession(ctx)

	// If no columns after restore, use default columns
	if len(e.state.Columns) == 0 {
		e.state.Columns = e.createDefaultColumns()
	}

	// Start auto-save timer
	e.startAutoSave()
}

// OnDismount performs cleanup when the component is unmounted.
// Stops the auto-save timer and performs a final save if there are changes.
func (e *Editor) OnDismount() {
	// Stop auto-save timer
	if !e.autoSaveTimer.IsUndefined() && !e.autoSaveTimer.IsNull() {
		app.Window().Call("clearInterval", e.autoSaveTimer)
	}

	// Final save if there are changes
	if e.state.HasChanges {
		e.saveState(context.Background())
	}
}

// Render renders the complete Excel-like editor interface.
// Creates a responsive layout with menu bar, main grid, and status bar components.
func (e *Editor) Render() app.UI {
	return app.Div().
		Class("excel-editor").
		TabIndex(0). // Make focusable for keyboard events
		OnKeyDown(e.handleKeyboardNavigation).
		Body(
			// Excel-like menu bar
			app.Header().
				Class("excel-menubar").
				Body(
					e.renderMenuBar(),
				),

			// Hidden file upload modal
			e.renderFileUploadModal(),

			// Main content area - always show grid
			app.Main().
				Class("excel-main").
				Body(
					e.renderDataGrid(),
				),

			// Status bar
			app.Footer().
				Class("excel-statusbar").
				Body(
					&StatusBar{
						SelectedCell: e.state.SelectedCell,
						CellValue:    e.getCellValue(e.state.SelectedCell),
						CellType:     e.getCellType(e.state.SelectedCell),
						TotalRows:    e.getDataRowCount(),
						TotalCols:    len(e.state.Columns),
						LastSaved:    e.state.LastSaved,
						HasChanges:   e.state.HasChanges,
						Error:        e.state.Error,
					},
				),
		)
}

// renderMenuBar renders the Excel-like menu bar with all functionality.
func (e *Editor) renderMenuBar() app.UI {
	return app.Div().
		Class("menubar-container").
		Body(
			// File menu
			app.Div().
				Class("menu-group").
				Body(
					app.Span().Class("menu-label").Text("File"),
					app.Button().
						Class("menu-button").
						Text("📄 New").
						Title("Create new spreadsheet (Ctrl+N)").
						OnClick(func(ctx app.Context, event app.Event) {
							e.handleNewFile()
							ctx.Update()
						}),
					app.Button().
						Class("menu-button").
						Text("📂 Open").
						Title("Open JSON schema file (Ctrl+O)").
						OnClick(func(ctx app.Context, event app.Event) {
							e.handleImportClick()
							ctx.Update()
						}),
					app.Button().
						Class("menu-button").
						Text("💾 Save").
						Title("Save current data (Ctrl+S)").
						OnClick(func(ctx app.Context, event app.Event) {
							e.handleSave()
							ctx.Update()
						}),
					app.Button().
						Class("menu-button").
						Text("📤 Export").
						Title("Export as JSON").
						OnClick(func(ctx app.Context, event app.Event) {
							e.handleExport()
							ctx.Update()
						}),
				),
			// Edit menu
			app.Div().
				Class("menu-group").
				Body(
					app.Span().Class("menu-label").Text("Edit"),
					app.Button().
						Class("menu-button").
						Text("🗑️ Clear").
						Title("Clear all data").
						OnClick(func(ctx app.Context, event app.Event) {
							e.handleClear()
							ctx.Update()
						}),
					app.Button().
						Class("menu-button").
						Text("➕ Add Row").
						Title("Add new row").
						OnClick(func(ctx app.Context, event app.Event) {
							e.handleAddRow()
							ctx.Update()
						}),
				),
			// Status indicator
			app.Div().
				Class("menu-status").
				Body(
					app.If(e.state.HasChanges,
						func() app.UI {
							return app.Span().Class("status-indicator unsaved").Text("● Unsaved changes")
						},
					).Else(
						func() app.UI {
							return app.Span().Class("status-indicator saved").Text("✓ Saved")
						},
					),
				),
		)
}

// renderFileUploadModal renders the file upload modal when triggered.
func (e *Editor) renderFileUploadModal() app.UI {
	if !e.showFileUpload {
		return app.Text("")
	}

	return app.Div().
		Class("modal-overlay").
		OnClick(func(ctx app.Context, event app.Event) {
			e.showFileUpload = false
			ctx.Update()
		}).
		Body(
			app.Div().
				Class("modal-content").
				OnClick(func(ctx app.Context, event app.Event) {
					event.Call("stopPropagation")
				}).
				Body(
					app.Div().
						Class("modal-header").
						Body(
							app.H2().Text("Upload JSON Schema"),
							app.Button().
								Class("modal-close").
								Text("✕").
								OnClick(func(ctx app.Context, event app.Event) {
									e.showFileUpload = false
									ctx.Update()
								}),
						),
					app.Div().
						Class("modal-body").
						Body(
							&FileUploadComponent{
								OnSchemaUploaded: e.handleSchemaUpload,
								OnError:          e.setError,
								AcceptTypes:      ".json,application/json",
								MaxFileSize:      1024 * 1024, // 1MB
							},
						),
				),
		)
}

// renderDataGrid renders the main data grid using the virtualized SpreadsheetGrid component.
func (e *Editor) renderDataGrid() app.UI {
	// Set up virtual canvas dimensions for Excel-like infinite scrolling
	virtualRows := 1000000 // 1M rows like Excel
	virtualCols := 1000    // 1K columns for performance
	
	// Use default columns if no schema is loaded
	cols := len(e.state.Columns)
	if cols == 0 {
		virtualCols = 26 // Default to A-Z columns
	} else {
		virtualCols = cols + 10 // Add extra columns beyond schema
	}

	return app.Div().
		Class("grid-container").
		Body(
			&grid.SpreadsheetGrid{
				VirtualRows: virtualRows,
				VirtualCols: virtualCols,
				Data:        e.state.Data,
				Size:        "medium",
				Scrollable:  true,
				MaxHeight:   "calc(100vh - 120px)", // Full height minus menu and status
				MaxWidth:    "100%",
				CellWidth:   120,  // Fixed width for performance
				CellHeight:  32,   // Fixed height for performance
				BufferSize:  5,    // Buffer for smooth scrolling
			},
		)
}

// Note: Cell editing is now handled directly by the virtualized SpreadsheetGrid component.
// The grid supports inline editing with double-click activation and handles validation internally.
// This provides better performance and user experience compared to overlay-based editing.

// Event Handlers

// handleKeyboardNavigation handles keyboard navigation within the editor.
// Implements Excel-like keyboard shortcuts and cell navigation.
func (e *Editor) handleKeyboardNavigation(ctx app.Context, event app.Event) {
	key := event.Get("key").String()
	ctrlKey := event.Get("ctrlKey").Bool()
	shiftKey := event.Get("shiftKey").Bool()

	// Handle keyboard shortcuts
	if ctrlKey {
		switch key {
		case "s":
			event.Call("preventDefault")
			e.handleSave()
			ctx.Update()
			return
		case "n":
			event.Call("preventDefault")
			e.handleNewFile()
			ctx.Update()
			return
		case "o":
			event.Call("preventDefault")
			e.handleImportClick()
			ctx.Update()
			return
		}
	}

	// Handle cell navigation
	if len(e.state.Columns) == 0 {
		return
	}

	parts := strings.Split(e.state.SelectedCell, "-")
	if len(parts) != 2 {
		return
	}

	row, _ := strconv.Atoi(parts[0])
	col, _ := strconv.Atoi(parts[1])
	maxRows := e.getDataRowCount()
	maxCols := len(e.state.Columns)

	switch key {
	case "ArrowUp":
		if row > 0 {
			e.selectCell(fmt.Sprintf("%d-%d", row-1, col))
		}
		event.Call("preventDefault")
	case "ArrowDown":
		if row < maxRows-1 {
			e.selectCell(fmt.Sprintf("%d-%d", row+1, col))
		} else {
			// Create new row if at bottom
			e.selectCell(fmt.Sprintf("%d-%d", maxRows, col))
		}
		event.Call("preventDefault")
	case "ArrowLeft":
		if col > 0 {
			e.selectCell(fmt.Sprintf("%d-%d", row, col-1))
		}
		event.Call("preventDefault")
	case "ArrowRight":
		if col < maxCols-1 {
			e.selectCell(fmt.Sprintf("%d-%d", row, col+1))
		}
		event.Call("preventDefault")
	case "Tab":
		if shiftKey {
			// Previous cell
			if col > 0 {
				e.selectCell(fmt.Sprintf("%d-%d", row, col-1))
			} else if row > 0 {
				e.selectCell(fmt.Sprintf("%d-%d", row-1, maxCols-1))
			}
		} else {
			// Next cell
			if col < maxCols-1 {
				e.selectCell(fmt.Sprintf("%d-%d", row, col+1))
			} else {
				e.selectCell(fmt.Sprintf("%d-%d", row+1, 0))
			}
		}
		event.Call("preventDefault")
	case "Enter":
		if e.state.IsEditing {
			e.finishCellEdit()
		} else {
			e.startCellEdit()
		}
		event.Call("preventDefault")
	case "Escape":
		if e.state.IsEditing {
			e.cancelCellEdit()
		}
		event.Call("preventDefault")
	case "F2":
		if !e.state.IsEditing {
			e.startCellEdit()
		}
		event.Call("preventDefault")
	default:
		// Start editing if typing a character
		if len(key) == 1 && !ctrlKey && !e.state.IsEditing {
			e.state.EditValue = key
			e.startCellEdit()
		}
	}

	ctx.Update()
}

// handleSchemaUpload processes uploaded JSON schema files.
// Validates the schema format and initializes the editor with column metadata.
func (e *Editor) handleSchemaUpload(schema *utils.JSONSchema, columns []grid.ColumnMetadata, sampleData []map[string]interface{}) {
	// Initialize editor with schema
	e.state.Schema = schema
	e.state.Columns = columns
	e.state.Data = make(map[string]string)
	e.state.SelectedCell = "0-0"
	e.showFileUpload = false // Close modal
	e.state.Error = ""
	e.markAsChanged()

	// Load sample data if provided
	if len(sampleData) > 0 {
		e.loadSampleData(sampleData)
	}

	// Create a new spreadsheet
	tableName := "Untitled Table"
	if schema.Title != "" {
		tableName = schema.Title
	}

	e.state.Spreadsheet = &data.Spreadsheet{
		ID:          fmt.Sprintf("excel-editor-%d", time.Now().Unix()),
		Name:        tableName,
		Description: schema.Description,
		Sheets:      []data.Sheet{},
		Owner:       "user",
		Created:     time.Now(),
		Updated:     time.Now(),
	}

	app.Log("Schema loaded successfully with", len(e.state.Columns), "columns")
}

// handleNewFile creates a new empty editor session.
func (e *Editor) handleNewFile() {
	if e.state.HasChanges {
		if !app.Window().Call("confirm", "You have unsaved changes. Are you sure you want to create a new file?").Bool() {
			return
		}
	}

	e.state = EditorState{
		Schema:       nil,
		Data:         make(map[string]string),
		Columns:      e.createDefaultColumns(), // Always have default columns
		SelectedCell: "0-0",
		LastSaved:    time.Now(),
	}
	e.showFileUpload = false
	e.state.Error = ""
}

// handleSave saves the current editor state to localStorage.
func (e *Editor) handleSave() {
	if err := e.saveState(context.Background()); err != nil {
		e.setError(fmt.Sprintf("Failed to save: %v", err))
		return
	}

	e.state.HasChanges = false
	e.state.LastSaved = time.Now()
	app.Log("File saved successfully")
}

// handleExport exports the current data as JSON.
func (e *Editor) handleExport() {
	var exportData interface{}

	if e.state.Schema != nil {
		// Export as proper JSON Schema with examples
		schema := *e.state.Schema // Copy the schema
		schema.Examples = e.convertDataToRows()
		exportData = schema
	} else {
		// Export as legacy format for backward compatibility
		exportData = map[string]interface{}{
			"schema": utils.LegacySchema{
				TableName:   e.getTableName(),
				Columns:     e.state.Columns,
				SampleData:  e.convertDataToRows(),
				Description: "Exported from Excel-like Editor",
			},
			"metadata": map[string]interface{}{
				"exported_at": time.Now().Format(time.RFC3339),
				"total_rows":  e.getDataRowCount(),
				"total_cols":  len(e.state.Columns),
			},
		}
	}

	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		e.setError(fmt.Sprintf("Export failed: %v", err))
		return
	}

	// Create download
	blob := app.Window().Get("Blob").New([]interface{}{string(jsonData)}, map[string]interface{}{
		"type": "application/json",
	})
	url := app.Window().Get("URL").Call("createObjectURL", blob)

	a := app.Window().Get("document").Call("createElement", "a")
	a.Set("href", url)

	fileName := fmt.Sprintf("%s-export.json", e.getTableName())
	if e.state.Schema != nil {
		fileName = fmt.Sprintf("%s-schema.json", e.getTableName())
	}
	a.Set("download", fileName)
	a.Call("click")

	app.Window().Get("URL").Call("revokeObjectURL", url)
	app.Log("Data exported successfully")
}

// handleImportClick triggers the file import dialog.
func (e *Editor) handleImportClick() {
	e.showFileUpload = true
}

// handleClear clears all data while keeping the schema.
func (e *Editor) handleClear() {
	if !app.Window().Call("confirm", "Are you sure you want to clear all data? This cannot be undone.").Bool() {
		return
	}

	e.state.Data = make(map[string]string)
	e.state.SelectedCell = "0-0"
	e.state.IsEditing = false
	e.state.EditValue = ""
	e.markAsChanged()
	app.Log("Data cleared")
}

// handleAddRow adds a new empty row to the spreadsheet.
func (e *Editor) handleAddRow() {
	currentRows := e.getDataRowCount()
	newRowIndex := currentRows

	// Initialize empty cells for the new row
	for colIndex := range e.state.Columns {
		cellKey := fmt.Sprintf("%d-%d", newRowIndex, colIndex)
		e.state.Data[cellKey] = ""
	}

	// Select the first cell of the new row
	e.state.SelectedCell = fmt.Sprintf("%d-0", newRowIndex)
	e.markAsChanged()
	app.Log("Added row", newRowIndex+1)
}

// handleCellInputChange handles changes to the cell edit input.
func (e *Editor) handleCellInputChange(value string) {
	e.state.EditValue = value
}

// handleCellSelect handles cell selection from the enhanced grid.
func (e *Editor) handleCellSelect(cellKey string) {
	if e.state.IsEditing {
		e.finishCellEdit()
	}
	e.state.SelectedCell = cellKey
}

// handleCellEditStart handles the start of cell editing from the enhanced grid.
func (e *Editor) handleCellEditStart(cellKey string) {
	e.state.SelectedCell = cellKey
	e.startCellEdit()
}

// handleCellValueChange handles cell value changes from the virtualized grid.
func (e *Editor) handleCellValueChange(cellKey, value string) error {
	// Validate the value first
	if err := e.validateCellValue(cellKey, value); err != nil {
		e.setError(fmt.Sprintf("Validation error: %v", err))
		return err
	}

	// Update the cell value
	e.state.Data[cellKey] = value
	e.state.Error = ""
	e.markAsChanged()
	return nil
}

// Cell Management Methods

// selectCell changes the currently selected cell.
func (e *Editor) selectCell(cellKey string) {
	if e.state.IsEditing {
		e.finishCellEdit()
	}
	e.state.SelectedCell = cellKey
}

// startCellEdit begins editing the currently selected cell.
func (e *Editor) startCellEdit() {
	e.state.IsEditing = true
	e.state.EditValue = e.getCellValue(e.state.SelectedCell)
}

// finishCellEdit completes cell editing and validates the input.
func (e *Editor) finishCellEdit() {
	if !e.state.IsEditing {
		return
	}

	// Validate the value
	if err := e.validateCellValue(e.state.SelectedCell, e.state.EditValue); err != nil {
		e.setError(fmt.Sprintf("Validation error: %v", err))
		return
	}

	// Update the cell value
	e.state.Data[e.state.SelectedCell] = e.state.EditValue
	e.state.IsEditing = false
	e.state.EditValue = ""
	e.state.Error = ""
	e.markAsChanged()
}

// cancelCellEdit cancels cell editing without saving changes.
func (e *Editor) cancelCellEdit() {
	e.state.IsEditing = false
	e.state.EditValue = ""
	e.state.Error = ""
}

// Helper Methods

// createDefaultColumns creates default columns for empty spreadsheet.
func (e *Editor) createDefaultColumns() []grid.ColumnMetadata {
	return []grid.ColumnMetadata{
		{
			Name:        "Column A",
			Type:        "varchar",
			Length:      255,
			Nullable:    true,
			Default:     "",
			Description: "Default column A",
		},
		{
			Name:        "Column B",
			Type:        "varchar",
			Length:      255,
			Nullable:    true,
			Default:     "",
			Description: "Default column B",
		},
		{
			Name:        "Column C",
			Type:        "varchar",
			Length:      255,
			Nullable:    true,
			Default:     "",
			Description: "Default column C",
		},
		{
			Name:        "Column D",
			Type:        "varchar",
			Length:      255,
			Nullable:    true,
			Default:     "",
			Description: "Default column D",
		},
		{
			Name:        "Column E",
			Type:        "varchar",
			Length:      255,
			Nullable:    true,
			Default:     "",
			Description: "Default column E",
		},
	}
}

// getCellValue returns the value of a cell or empty string if not found.
func (e *Editor) getCellValue(cellKey string) string {
	return e.state.Data[cellKey]
}

// getCellType returns the data type of a cell based on its column.
func (e *Editor) getCellType(cellKey string) string {
	parts := strings.Split(cellKey, "-")
	if len(parts) != 2 {
		return "text"
	}

	col, err := strconv.Atoi(parts[1])
	if err != nil || col >= len(e.state.Columns) {
		return "text"
	}

	return e.state.Columns[col].Type
}

// getTableName returns the current table name or a default.
func (e *Editor) getTableName() string {
	if e.state.Spreadsheet != nil && e.state.Spreadsheet.Name != "" {
		return e.state.Spreadsheet.Name
	}
	return "Untitled Table"
}

// getDataRowCount calculates the number of data rows.
func (e *Editor) getDataRowCount() int {
	maxRow := -1
	for key := range e.state.Data {
		parts := strings.Split(key, "-")
		if len(parts) == 2 {
			if row, err := strconv.Atoi(parts[0]); err == nil && row > maxRow {
				maxRow = row
			}
		}
	}
	return maxRow + 1
}

// validateCellValue validates a cell value against its column type and JSON Schema constraints.
func (e *Editor) validateCellValue(cellKey, value string) error {
	if value == "" {
		// Check if this column is required
		parts := strings.Split(cellKey, "-")
		if len(parts) == 2 {
			if col, err := strconv.Atoi(parts[1]); err == nil && col < len(e.state.Columns) {
				column := e.state.Columns[col]
				if !column.Nullable {
					return fmt.Errorf("this field is required")
				}
			}
		}
		return nil // Empty values are allowed for nullable fields
	}

	// Get column information for enhanced validation
	parts := strings.Split(cellKey, "-")
	if len(parts) != 2 {
		return nil
	}

	col, err := strconv.Atoi(parts[1])
	if err != nil || col >= len(e.state.Columns) {
		return nil
	}

	column := e.state.Columns[col]
	cellType := column.Type

	// Basic type validation
	switch cellType {
	case "integer", "bigint":
		if val, err := strconv.ParseInt(value, 10, 64); err != nil {
			return fmt.Errorf("must be a valid integer")
		} else {
			// Additional JSON Schema validation if available
			if e.state.Schema != nil && e.state.Schema.Properties != nil {
				if prop, exists := e.state.Schema.Properties[column.Name]; exists {
					if prop.Minimum != nil && float64(val) < *prop.Minimum {
						return fmt.Errorf("must be at least %v", *prop.Minimum)
					}
					if prop.Maximum != nil && float64(val) > *prop.Maximum {
						return fmt.Errorf("must be at most %v", *prop.Maximum)
					}
				}
			}
		}
	case "decimal", "numeric":
		if val, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("must be a valid number")
		} else {
			// Additional JSON Schema validation if available
			if e.state.Schema != nil && e.state.Schema.Properties != nil {
				if prop, exists := e.state.Schema.Properties[column.Name]; exists {
					if prop.Minimum != nil && val < *prop.Minimum {
						return fmt.Errorf("must be at least %v", *prop.Minimum)
					}
					if prop.Maximum != nil && val > *prop.Maximum {
						return fmt.Errorf("must be at most %v", *prop.Maximum)
					}
				}
			}
		}
	case "boolean":
		lower := strings.ToLower(value)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return fmt.Errorf("must be true, false, 1, or 0")
		}
	case "date":
		if _, err := time.Parse("2006-01-02", value); err != nil {
			return fmt.Errorf("must be a valid date (YYYY-MM-DD)")
		}
	case "timestamp":
		formats := []string{
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
		}
		var parseErr error
		valid := false
		for _, format := range formats {
			if _, parseErr = time.Parse(format, value); parseErr == nil {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("must be a valid timestamp (e.g., 2023-01-01T10:30:00)")
		}
	case "varchar", "text":
		// String length validation
		if column.Length > 0 && len(value) > column.Length {
			return fmt.Errorf("must be at most %d characters", column.Length)
		}

		// Additional JSON Schema validation if available
		if e.state.Schema != nil && e.state.Schema.Properties != nil {
			if prop, exists := e.state.Schema.Properties[column.Name]; exists {
				if prop.MinLength != nil && len(value) < *prop.MinLength {
					return fmt.Errorf("must be at least %d characters", *prop.MinLength)
				}
				if prop.MaxLength != nil && len(value) > *prop.MaxLength {
					return fmt.Errorf("must be at most %d characters", *prop.MaxLength)
				}

				// Format validation for strings
				if prop.Format != "" {
					if err := e.validateStringFormat(value, prop.Format); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// validateStringFormat validates string values against JSON Schema format constraints
func (e *Editor) validateStringFormat(value, format string) error {
	switch format {
	case "email":
		if !utils.IsValidEmail(value) {
			return fmt.Errorf("must be a valid email address")
		}
	case "uri", "url":
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			return fmt.Errorf("must be a valid URL")
		}
	case "date":
		if _, err := time.Parse("2006-01-02", value); err != nil {
			return fmt.Errorf("must be a valid date (YYYY-MM-DD)")
		}
	case "date-time":
		formats := []string{
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05",
		}
		valid := false
		for _, timeFormat := range formats {
			if _, err := time.Parse(timeFormat, value); err == nil {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("must be a valid date-time (e.g., 2023-01-01T10:30:00)")
		}
	case "time":
		if _, err := time.Parse("15:04:05", value); err != nil {
			return fmt.Errorf("must be a valid time (HH:MM:SS)")
		}
	}
	return nil
}

// loadSampleData loads sample data from the schema into the grid.
func (e *Editor) loadSampleData(sampleData []map[string]interface{}) {
	for rowIndex, rowData := range sampleData {
		for colIndex, column := range e.state.Columns {
			if value, exists := rowData[column.Name]; exists {
				cellKey := fmt.Sprintf("%d-%d", rowIndex, colIndex)
				e.state.Data[cellKey] = fmt.Sprintf("%v", value)
			}
		}
	}
}

// convertDataToRows converts the internal data format to row-based format for export.
func (e *Editor) convertDataToRows() []map[string]interface{} {
	rows := make([]map[string]interface{}, e.getDataRowCount())

	for i := range rows {
		rows[i] = make(map[string]interface{})
	}

	for cellKey, value := range e.state.Data {
		parts := strings.Split(cellKey, "-")
		if len(parts) == 2 {
			row, _ := strconv.Atoi(parts[0])
			col, _ := strconv.Atoi(parts[1])

			if row < len(rows) && col < len(e.state.Columns) {
				columnName := e.state.Columns[col].Name
				rows[row][columnName] = value
			}
		}
	}

	return rows
}

// markAsChanged marks the editor state as having unsaved changes.
func (e *Editor) markAsChanged() {
	e.state.HasChanges = true
	e.state.Spreadsheet.Updated = time.Now()
}

// setError sets an error message in the editor state.
func (e *Editor) setError(message string) {
	e.state.Error = message
	app.Log("Error:", message)
}

// Storage Methods

// saveState saves the current editor state to localStorage.
func (e *Editor) saveState(ctx context.Context) error {
	if e.storage == nil {
		return fmt.Errorf("storage not initialized")
	}

	// Save the spreadsheet if it exists
	if e.state.Spreadsheet != nil {
		if err := e.storage.SaveSpreadsheet(ctx, e.state.Spreadsheet); err != nil {
			return err
		}
	}

	// Save editor state
	stateData := map[string]interface{}{
		"columns":       e.state.Columns,
		"data":          e.state.Data,
		"selected_cell": e.state.SelectedCell,
		"last_saved":    time.Now().Unix(),
	}

	data, err := json.Marshal(stateData)
	if err != nil {
		return err
	}

	app.Window().Get("localStorage").Call("setItem", e.StoragePrefix+"_editor_state", string(data))
	return nil
}

// restoreLastSession attempts to restore the last editor session from localStorage.
func (e *Editor) restoreLastSession(ctx app.Context) {
	stateKey := e.StoragePrefix + "_editor_state"
	value := app.Window().Get("localStorage").Call("getItem", stateKey)

	if value.IsNull() || value.IsUndefined() {
		return
	}

	var stateData map[string]interface{}
	if err := json.Unmarshal([]byte(value.String()), &stateData); err != nil {
		app.Log("Failed to restore session:", err)
		return
	}

	// Restore columns
	if columnsData, ok := stateData["columns"].([]interface{}); ok {
		for _, columnInterface := range columnsData {
			if columnMap, ok := columnInterface.(map[string]interface{}); ok {
				column := grid.ColumnMetadata{
					Name:        getString(columnMap, "name"),
					Type:        getString(columnMap, "type"),
					Nullable:    getBool(columnMap, "nullable"),
					Default:     getString(columnMap, "default"),
					Description: getString(columnMap, "description"),
					Length:      getInt(columnMap, "length"),
					Precision:   getInt(columnMap, "precision"),
					Scale:       getInt(columnMap, "scale"),
				}
				e.state.Columns = append(e.state.Columns, column)
			}
		}
	}

	// Restore data
	if dataMap, ok := stateData["data"].(map[string]interface{}); ok {
		for key, value := range dataMap {
			if strValue, ok := value.(string); ok {
				e.state.Data[key] = strValue
			}
		}
	}

	// Restore selected cell
	if selectedCell, ok := stateData["selected_cell"].(string); ok {
		e.state.SelectedCell = selectedCell
	}

	if len(e.state.Columns) > 0 {
		app.Log("Session restored with", len(e.state.Columns), "columns")
	}
}

// startAutoSave starts the auto-save timer.
func (e *Editor) startAutoSave() {
	interval := int(e.AutoSaveInterval.Milliseconds())
	e.autoSaveTimer = app.Window().Call("setInterval", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		if e.state.HasChanges {
			if err := e.saveState(context.Background()); err != nil {
				app.Log("Auto-save failed:", err)
			} else {
				e.state.HasChanges = false
				e.state.LastSaved = time.Now()
				app.Log("Auto-saved at", time.Now().Format("15:04:05"))
			}
		}
		return nil
	}), interval)
}

// Utility functions for type conversion
func getString(m map[string]interface{}, key string) string {
	if value, ok := m[key].(string); ok {
		return value
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if value, ok := m[key].(bool); ok {
		return value
	}
	return false
}

func getInt(m map[string]interface{}, key string) int {
	if value, ok := m[key].(float64); ok {
		return int(value)
	}
	return 0
}

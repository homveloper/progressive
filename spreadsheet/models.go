package spreadsheet

import (
	"encoding/json"
	"fmt"
	"time"
)

// =============================================================================
// Workbook Data Model - Inspired by Univer's architecture
// =============================================================================

// WorkbookData represents the complete state of a spreadsheet workbook.
// It follows Univer's architecture principles with value-based semantics for
// immutability and thread safety. The workbook contains multiple sheets,
// shared styles, and metadata.
//
// Key features:
//   - Immutable value semantics (no pointers except for zero-copy scenarios)
//   - Thread-safe operations through value copying
//   - Hierarchical structure: Workbook -> Sheets -> Cells
//   - Shared style definitions across sheets
//   - Comprehensive metadata tracking
//
// Example usage:
//
//	workbook := NewWorkbook("Financial Report 2024")
//	workbook = workbook.UpdateCell("sheet1", 0, 0, CellValue{
//	    Value: "Revenue",
//	    Type:  CellTypeString,
//	})
type WorkbookData struct {
	// ID is a unique identifier for the workbook
	ID string `json:"id"`
	
	// Name is the human-readable name of the workbook
	Name string `json:"name"`
	
	// AppVersion tracks the application version that created/modified this workbook
	AppVersion string `json:"appVersion"`
	
	// Locale specifies the locale for formatting and display (e.g., "en-US", "ko-KR")
	Locale string `json:"locale"`
	
	// CreatedAt records when the workbook was initially created
	CreatedAt time.Time `json:"createdAt"`
	
	// UpdatedAt tracks the last modification time
	UpdatedAt time.Time `json:"updatedAt"`
	
	// Styles contains shared style definitions used across all sheets
	Styles map[string]StyleData `json:"styles"`
	
	// SheetOrder defines the order of sheet tabs as they appear in the UI
	SheetOrder []string `json:"sheetOrder"`
	
	// Sheets contains all worksheets indexed by their ID
	Sheets map[string]SheetData `json:"sheets"`
	
	// Resources stores additional workbook-level resources (images, charts, etc.)
	Resources map[string]any `json:"resources,omitempty"`
}

// =============================================================================
// Worksheet Data Model
// =============================================================================

// SheetData represents a single worksheet within a workbook
type SheetData struct {
	ID                string                    `json:"id"`
	Name              string                    `json:"name"`
	TabColor          string                    `json:"tabColor,omitempty"`
	Hidden            bool                      `json:"hidden"`
	Protected         bool                      `json:"protected"`
	RowCount          int                       `json:"rowCount"`
	ColumnCount       int                       `json:"columnCount"`
	DefaultRowHeight  float64                   `json:"defaultRowHeight"`
	DefaultColWidth   float64                   `json:"defaultColumnWidth"`
	Freeze            FreezeConfig              `json:"freeze,omitempty"`
	MergeData         []MergeRange              `json:"mergeData,omitempty"`
	CellData          map[int]map[int]CellValue `json:"cellData"`
	RowData           map[int]RowConfig         `json:"rowData,omitempty"`
	ColumnData        map[int]ColumnConfig      `json:"columnData,omitempty"`
	ConditionalFormat []ConditionalFormat       `json:"conditionalFormat,omitempty"`
	ShowGridlines     bool                      `json:"showGridlines"`
	RightToLeft       bool                      `json:"rightToLeft"`
	DefaultStyle      StyleData                 `json:"defaultStyle,omitempty"`
	UpdatedAt         time.Time                 `json:"updatedAt"`
}

// FreezeConfig represents frozen panes configuration
type FreezeConfig struct {
	RowCount    int `json:"rowCount"`
	ColumnCount int `json:"columnCount"`
}

// MergeRange represents merged cell range
type MergeRange struct {
	StartRow int `json:"startRow"`
	StartCol int `json:"startCol"`
	EndRow   int `json:"endRow"`
	EndCol   int `json:"endCol"`
}

// RowConfig represents row-specific configuration
type RowConfig struct {
	Height float64 `json:"height,omitempty"`
	Hidden bool    `json:"hidden"`
	Style  string  `json:"style,omitempty"`
}

// ColumnConfig represents column-specific configuration
type ColumnConfig struct {
	Width  float64 `json:"width,omitempty"`
	Hidden bool    `json:"hidden"`
	Style  string  `json:"style,omitempty"`
}

// =============================================================================
// Cell Data Model
// =============================================================================

// CellType enumeration for different cell content types
type CellType int

const (
	CellTypeString  CellType = 1
	CellTypeNumber  CellType = 2
	CellTypeBoolean CellType = 3
	CellTypeForceText CellType = 4
	CellTypeError   CellType = 5
)

// CellValue represents the complete data for a single cell
type CellValue struct {
	// Core value - can be string, number, or boolean
	Value interface{} `json:"v,omitempty"`
	
	// Cell type
	Type CellType `json:"t,omitempty"`
	
	// Style can be either a style ID (string) or StyleData
	Style interface{} `json:"s,omitempty"` // Keep as interface for flexibility
	
	// Formula string
	Formula string `json:"f,omitempty"`
	
	// Formula ID for shared formulas
	FormulaID string `json:"si,omitempty"`
	
	// Rich text content
	RichText RichTextData `json:"p,omitempty"`
	
	// Cell metadata
	Metadata CellMetadata `json:"m,omitempty"`
	
	// Custom fields for extensions
	Custom map[string]interface{} `json:"custom,omitempty"`
}

// CellMetadata stores additional cell information
type CellMetadata struct {
	Comment      string    `json:"comment,omitempty"`
	LastModified time.Time `json:"lastModified,omitempty"`
	ModifiedBy   string    `json:"modifiedBy,omitempty"`
	Locked       bool      `json:"locked"`
	Validation   DataValidation `json:"validation,omitempty"`
}

// DataValidation represents cell validation rules
type DataValidation struct {
	Type         string      `json:"type"` // list, number, date, text, custom
	Formula1     string      `json:"formula1,omitempty"`
	Formula2     string      `json:"formula2,omitempty"`
	AllowBlank   bool        `json:"allowBlank"`
	ShowDropdown bool        `json:"showDropdown"`
	ErrorTitle   string      `json:"errorTitle,omitempty"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
	Options      []string    `json:"options,omitempty"`
}

// =============================================================================
// Style Data Model
// =============================================================================

// StyleData represents comprehensive cell styling
type StyleData struct {
	Font       FontStyle       `json:"font,omitempty"`
	Background BackgroundStyle `json:"background,omitempty"`
	Border     BorderStyle     `json:"border,omitempty"`
	Alignment  AlignmentStyle  `json:"alignment,omitempty"`
	Format     NumberFormat    `json:"format,omitempty"`
}

// FontStyle represents text formatting
type FontStyle struct {
	Name      string  `json:"name,omitempty"`
	Size      float64 `json:"size,omitempty"`
	Bold      bool    `json:"bold"`
	Italic    bool    `json:"italic"`
	Underline string  `json:"underline,omitempty"` // none, single, double
	Strike    bool    `json:"strike"`
	Color     string  `json:"color,omitempty"`
}

// BackgroundStyle represents cell background
type BackgroundStyle struct {
	Color   string `json:"color,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

// BorderStyle represents cell borders
type BorderStyle struct {
	Top    BorderLine `json:"top,omitempty"`
	Right  BorderLine `json:"right,omitempty"`
	Bottom BorderLine `json:"bottom,omitempty"`
	Left   BorderLine `json:"left,omitempty"`
}

// BorderLine represents a single border line
type BorderLine struct {
	Style string `json:"style"` // thin, medium, thick, dashed, dotted
	Color string `json:"color"`
}

// AlignmentStyle represents text alignment
type AlignmentStyle struct {
	Horizontal   string `json:"horizontal,omitempty"` // left, center, right, justify
	Vertical     string `json:"vertical,omitempty"`   // top, middle, bottom
	WrapText     bool   `json:"wrapText"`
	Indent       int    `json:"indent,omitempty"`
	TextRotation int    `json:"textRotation,omitempty"` // -90 to 90
}

// NumberFormat represents number formatting
type NumberFormat struct {
	Pattern  string `json:"pattern"`  // e.g., "#,##0.00", "$#,##0.00"
	Type     string `json:"type"`     // number, currency, percentage, date, time
	Currency string `json:"currency,omitempty"`
}

// =============================================================================
// Rich Text Data Model
// =============================================================================

// RichTextData represents formatted text with multiple runs
type RichTextData struct {
	Body DocumentBody `json:"body,omitempty"`
}

// DocumentBody contains paragraphs of rich text
type DocumentBody struct {
	Paragraphs []Paragraph `json:"paragraphs"`
}

// Paragraph represents a paragraph in rich text
type Paragraph struct {
	Elements []TextRun `json:"elements"`
	Style    ParagraphStyle `json:"style,omitempty"`
}

// TextRun represents a run of text with consistent formatting
type TextRun struct {
	Text  string     `json:"text"`
	Style FontStyle `json:"style,omitempty"`
}

// ParagraphStyle represents paragraph-level styling
type ParagraphStyle struct {
	SpaceBefore float64 `json:"spaceBefore,omitempty"`
	SpaceAfter  float64 `json:"spaceAfter,omitempty"`
	LineSpacing float64 `json:"lineSpacing,omitempty"`
}

// =============================================================================
// Conditional Formatting
// =============================================================================

// ConditionalFormat represents conditional formatting rules
type ConditionalFormat struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"` // cellValue, colorScale, dataBar, iconSet
	Priority int         `json:"priority"`
	Range    CellRange  `json:"range"`
	Rule     interface{} `json:"rule"`
	Style    StyleData  `json:"style,omitempty"`
}

// CellRange represents a range of cells
type CellRange struct {
	StartRow int `json:"startRow"`
	StartCol int `json:"startCol"`
	EndRow   int `json:"endRow"`
	EndCol   int `json:"endCol"`
}

// =============================================================================
// Helper Methods
// =============================================================================

// NewWorkbook creates a new workbook with default settings
func NewWorkbook(name string) WorkbookData {
	defaultSheet := NewSheet("Sheet1")
	
	workbook := WorkbookData{
		ID:         GenerateID(),
		Name:       name,
		AppVersion: "1.0.0",
		Locale:     "en-US",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Styles:     make(map[string]StyleData),
		SheetOrder: []string{defaultSheet.ID},
		Sheets:     map[string]SheetData{defaultSheet.ID: defaultSheet},
		Resources:  make(map[string]any),
	}
	
	return workbook
}

// NewSheet creates a new worksheet with default settings
func NewSheet(name string) SheetData {
	return SheetData{
		ID:               GenerateID(),
		Name:             name,
		RowCount:         100,
		ColumnCount:      26,
		DefaultRowHeight: 24,
		DefaultColWidth:  64,
		CellData:         make(map[int]map[int]CellValue),
		ShowGridlines:    true,
		RowData:          make(map[int]RowConfig),
		ColumnData:       make(map[int]ColumnConfig),
		UpdatedAt:        time.Now(),
	}
}

// GetCell retrieves a cell from the sheet
func (s SheetData) GetCell(row, col int) (CellValue, bool) {
	if rowData, ok := s.CellData[row]; ok {
		if cell, ok := rowData[col]; ok {
			return cell, true
		}
	}
	return CellValue{}, false
}

// SetCell sets a cell value in the sheet (returns new SheetData)
func (s SheetData) SetCell(row, col int, cell CellValue) SheetData {
	if s.CellData == nil {
		s.CellData = make(map[int]map[int]CellValue)
	}
	if s.CellData[row] == nil {
		s.CellData[row] = make(map[int]CellValue)
	}
	s.CellData[row][col] = cell
	s.UpdatedAt = time.Now()
	return s
}

// ToJSON serializes the workbook to JSON
func (w WorkbookData) ToJSON() ([]byte, error) {
	return json.MarshalIndent(w, "", "  ")
}

// FromJSON deserializes a workbook from JSON
func FromJSON(data []byte) (WorkbookData, error) {
	var workbook WorkbookData
	err := json.Unmarshal(data, &workbook)
	return workbook, err
}

// GenerateID generates a unique identifier
func GenerateID() string {
	// Simple implementation - in production, use UUID
	return fmt.Sprintf("id_%d", time.Now().UnixNano())
}

// =============================================================================
// Helper Methods for Value-Based Updates
// =============================================================================

// UpdateSheet updates a sheet in the workbook and returns the new workbook
func (w WorkbookData) UpdateSheet(sheetID string, sheet SheetData) WorkbookData {
	if w.Sheets == nil {
		w.Sheets = make(map[string]SheetData)
	}
	w.Sheets[sheetID] = sheet
	w.UpdatedAt = time.Now()
	return w
}

// GetSheet retrieves a sheet from the workbook
func (w WorkbookData) GetSheet(sheetID string) (SheetData, bool) {
	sheet, ok := w.Sheets[sheetID]
	return sheet, ok
}

// AddSheet adds a new sheet to the workbook
func (w WorkbookData) AddSheet(sheet SheetData) WorkbookData {
	if w.Sheets == nil {
		w.Sheets = make(map[string]SheetData)
	}
	w.Sheets[sheet.ID] = sheet
	w.SheetOrder = append(w.SheetOrder, sheet.ID)
	w.UpdatedAt = time.Now()
	return w
}

// UpdateCell updates a cell in a specific sheet
func (w WorkbookData) UpdateCell(sheetID string, row, col int, cell CellValue) WorkbookData {
	if sheet, ok := w.Sheets[sheetID]; ok {
		updatedSheet := sheet.SetCell(row, col, cell)
		w.Sheets[sheetID] = updatedSheet
		w.UpdatedAt = time.Now()
	}
	return w
}

// GetCell retrieves a cell from a specific sheet
func (w WorkbookData) GetCell(sheetID string, row, col int) (CellValue, bool) {
	if sheet, ok := w.Sheets[sheetID]; ok {
		return sheet.GetCell(row, col)
	}
	return CellValue{}, false
}
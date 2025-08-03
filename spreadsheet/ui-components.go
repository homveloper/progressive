package spreadsheet

import (
	"fmt"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// =============================================================================
// Basic UI Components Library
// =============================================================================

// ButtonVariant represents the visual style variant of a button.
// Each variant corresponds to different semantic meanings and color schemes.
type ButtonVariant string

const (
	// ButtonPrimary represents the main action button (blue)
	ButtonPrimary ButtonVariant = "primary"
	
	// ButtonSecondary represents secondary actions (gray)
	ButtonSecondary ButtonVariant = "secondary"
	
	// ButtonSuccess represents success or positive actions (green)
	ButtonSuccess ButtonVariant = "success"
	
	// ButtonDanger represents destructive or dangerous actions (red)
	ButtonDanger ButtonVariant = "danger"
	
	// ButtonWarning represents warning or caution actions (yellow)
	ButtonWarning ButtonVariant = "warning"
	
	// ButtonInfo represents informational actions (light blue)
	ButtonInfo ButtonVariant = "info"
	
	// ButtonOutline represents a outlined version with transparent background
	ButtonOutline ButtonVariant = "outline"
)

// ButtonSize represents the size variant of a button.
// Different sizes provide visual hierarchy and fit different UI contexts.
type ButtonSize string

const (
	// ButtonSmall creates a compact button suitable for tight spaces
	ButtonSmall ButtonSize = "sm"
	
	// ButtonMedium creates a standard-sized button for most use cases  
	ButtonMedium ButtonSize = "md"
	
	// ButtonLarge creates a prominent button for primary actions
	ButtonLarge ButtonSize = "lg"
)

// UIButton is a reusable button component with multiple visual variants and sizes.
// It supports icons, loading states, and disabled states for comprehensive UI needs.
//
// The component provides:
//   - Multiple visual variants (primary, secondary, success, danger, etc.)
//   - Three size options (small, medium, large)
//   - Optional icon display
//   - Loading state with spinner
//   - Disabled state
//   - Click event handling
//
// Example usage:
//
//	button := &UIButton{
//	    Text:     "Save Changes",
//	    Variant:  ButtonPrimary,
//	    Size:     ButtonMedium,
//	    Icon:     "💾",
//	    OnClick:  func() { fmt.Println("Save clicked!") },
//	}
type UIButton struct {
	app.Compo
	
	// Text is the button label text
	Text string
	
	// Variant determines the visual style and color scheme
	Variant ButtonVariant
	
	// Size determines the button dimensions and font size
	Size ButtonSize
	
	// Disabled prevents user interaction when true
	Disabled bool
	
	// Icon is an optional icon/emoji to display alongside text
	Icon string
	
	// Loading shows a spinner and disables interaction when true
	Loading bool
	
	// OnClick is the callback function executed when the button is clicked
	OnClick func()
}

func (b *UIButton) Render() app.UI {
	classes := []string{"ui-button"}
	classes = append(classes, fmt.Sprintf("btn-%s", b.Variant))
	classes = append(classes, fmt.Sprintf("btn-%s", b.Size))
	
	if b.Disabled || b.Loading {
		classes = append(classes, "disabled")
	}
	
	className := strings.Join(classes, " ")
	
	var content []app.UI
	
	if b.Loading {
		content = append(content, app.Span().Class("loading-spinner"))
	}
	
	if b.Icon != "" && !b.Loading {
		content = append(content, app.Span().Class("btn-icon").Text(b.Icon))
	}
	
	if b.Text != "" {
		content = append(content, app.Span().Class("btn-text").Text(b.Text))
	}
	
	return app.Button().
		Class(className).
		Disabled(b.Disabled || b.Loading).
		OnClick(func(ctx app.Context, e app.Event) {
			if b.OnClick != nil && !b.Disabled && !b.Loading {
				b.OnClick()
			}
		}).
		Body(content...)
}

// Badge variants
type BadgeVariant string

const (
	BadgePrimary   BadgeVariant = "primary"
	BadgeSecondary BadgeVariant = "secondary"
	BadgeSuccess   BadgeVariant = "success"
	BadgeDanger    BadgeVariant = "danger"
	BadgeWarning   BadgeVariant = "warning"
	BadgeInfo      BadgeVariant = "info"
)

// UIBadge - 배지 컴포넌트
type UIBadge struct {
	app.Compo
	Text    string
	Variant BadgeVariant
	Dot     bool
	Count   int
}

func (b *UIBadge) Render() app.UI {
	classes := []string{"ui-badge"}
	classes = append(classes, fmt.Sprintf("badge-%s", b.Variant))
	
	if b.Dot {
		classes = append(classes, "badge-dot")
	}
	
	className := strings.Join(classes, " ")
	
	text := b.Text
	if b.Count > 0 {
		if b.Count > 99 {
			text = "99+"
		} else {
			text = fmt.Sprintf("%d", b.Count)
		}
	}
	
	return app.Span().Class(className).Text(text)
}

// UICard - 카드 컴포넌트
type UICard struct {
	app.Compo
	Title       string
	Description string
	Image       string
	Content     app.UI
	Actions     []app.UI
	Hoverable   bool
	Bordered    bool
}

func (c *UICard) Render() app.UI {
	classes := []string{"ui-card"}
	
	if c.Hoverable {
		classes = append(classes, "card-hoverable")
	}
	
	if c.Bordered {
		classes = append(classes, "card-bordered")
	}
	
	className := strings.Join(classes, " ")
	
	var cardContent []app.UI
	
	// Image
	if c.Image != "" {
		cardContent = append(cardContent, 
			app.Div().Class("card-image").Body(
				app.Img().Src(c.Image).Alt(c.Title),
			),
		)
	}
	
	// Content area
	contentArea := []app.UI{}
	
	if c.Title != "" {
		contentArea = append(contentArea, 
			app.H3().Class("card-title").Text(c.Title),
		)
	}
	
	if c.Description != "" {
		contentArea = append(contentArea, 
			app.P().Class("card-description").Text(c.Description),
		)
	}
	
	if c.Content != nil {
		contentArea = append(contentArea, c.Content)
	}
	
	if len(contentArea) > 0 {
		cardContent = append(cardContent, 
			app.Div().Class("card-content").Body(contentArea...),
		)
	}
	
	// Actions
	if len(c.Actions) > 0 {
		cardContent = append(cardContent,
			app.Div().Class("card-actions").Body(c.Actions...),
		)
	}
	
	return app.Div().Class(className).Body(cardContent...)
}

// CellFormat represents cell formatting options
type CellFormat struct {
	Bold      bool
	Italic    bool
	Underline bool
	Color     string
	BgColor   string
}

// SpreadsheetCell - 스프레드시트 셀 컴포넌트
type SpreadsheetCell struct {
	app.Compo
	Value     string
	Selected  bool
	Editing   bool
	Format    CellFormat
	ReadOnly  bool
	OnClick   func()
	OnChange  func(value string)
	OnKeyDown func(key string)
}

func (c *SpreadsheetCell) Render() app.UI {
	classes := []string{"spreadsheet-cell"}
	
	if c.Selected {
		classes = append(classes, "cell-selected")
	}
	
	if c.ReadOnly {
		classes = append(classes, "cell-readonly")
	}
	
	className := strings.Join(classes, " ")
	
	var cellStyle strings.Builder
	if c.Format.Bold {
		cellStyle.WriteString("font-weight: bold; ")
	}
	if c.Format.Italic {
		cellStyle.WriteString("font-style: italic; ")
	}
	if c.Format.Underline {
		cellStyle.WriteString("text-decoration: underline; ")
	}
	if c.Format.Color != "" {
		cellStyle.WriteString(fmt.Sprintf("color: %s; ", c.Format.Color))
	}
	if c.Format.BgColor != "" {
		cellStyle.WriteString(fmt.Sprintf("background-color: %s; ", c.Format.BgColor))
	}
	
	if c.Editing {
		return app.Input().
			Class("cell-input").
			Type("text").
			Value(c.Value).
			AutoFocus(true).
			Style("style", cellStyle.String()).
			OnChange(func(ctx app.Context, e app.Event) {
				if c.OnChange != nil {
					c.OnChange(e.Get("target").Get("value").String())
				}
			}).
			OnKeyDown(func(ctx app.Context, e app.Event) {
				if c.OnKeyDown != nil {
					c.OnKeyDown(e.Get("key").String())
				}
			})
	}
	
	return app.Div().
		Class(className).
		Style("style", cellStyle.String()).
		Text(c.Value).
		OnClick(func(ctx app.Context, e app.Event) {
			if c.OnClick != nil && !c.ReadOnly {
				c.OnClick()
			}
		})
}

// SheetTab - 시트 탭 컴포넌트
type SheetTab struct {
	app.Compo
	Name     string
	Active   bool
	Closable bool
	OnClick  func()
	OnClose  func()
}

func (t *SheetTab) Render() app.UI {
	classes := []string{"sheet-tab"}
	
	if t.Active {
		classes = append(classes, "tab-active")
	}
	
	className := strings.Join(classes, " ")
	
	var tabContent []app.UI
	
	tabContent = append(tabContent, 
		app.Span().Class("tab-name").Text(t.Name),
	)
	
	if t.Closable {
		tabContent = append(tabContent,
			app.Button().Class("tab-close").Text("×").
				OnClick(func(ctx app.Context, e app.Event) {
					if t.OnClose != nil {
						t.OnClose()
					}
				}),
		)
	}
	
	return app.Div().
		Class(className).
		OnClick(func(ctx app.Context, e app.Event) {
			if t.OnClick != nil {
				t.OnClick()
			}
		}).
		Body(tabContent...)
}

// ToolbarButton - 툴바 버튼 컴포넌트
type ToolbarButton struct {
	app.Compo
	Icon     string
	Text     string
	Tooltip  string
	Active   bool
	Disabled bool
	OnClick  func()
}

func (t *ToolbarButton) Render() app.UI {
	classes := []string{"toolbar-button"}
	
	if t.Active {
		classes = append(classes, "btn-active")
	}
	
	if t.Disabled {
		classes = append(classes, "btn-disabled")
	}
	
	className := strings.Join(classes, " ")
	
	var content []app.UI
	
	if t.Icon != "" {
		content = append(content, app.Span().Class("btn-icon").Text(t.Icon))
	}
	
	if t.Text != "" {
		content = append(content, app.Span().Class("btn-text").Text(t.Text))
	}
	
	return app.Button().
		Class(className).
		Title(t.Tooltip).
		Disabled(t.Disabled).
		OnClick(func(ctx app.Context, e app.Event) {
			if t.OnClick != nil && !t.Disabled {
				t.OnClick()
			}
		}).
		Body(content...)
}

// InputField - 입력 필드 컴포넌트
type InputField struct {
	app.Compo
	Label       string
	Value       string
	Placeholder string
	Type        string
	Required    bool
	Disabled    bool
	Error       string
	OnChange    func(value string)
}

func (i *InputField) Render() app.UI {
	classes := []string{"input-field"}
	
	if i.Error != "" {
		classes = append(classes, "field-error")
	}
	
	if i.Disabled {
		classes = append(classes, "field-disabled")
	}
	
	className := strings.Join(classes, " ")
	
	var fieldContent []app.UI
	
	if i.Label != "" {
		labelClasses := []string{"field-label"}
		if i.Required {
			labelClasses = append(labelClasses, "label-required")
		}
		
		fieldContent = append(fieldContent,
			app.Label().Class(strings.Join(labelClasses, " ")).Text(i.Label),
		)
	}
	
	inputType := i.Type
	if inputType == "" {
		inputType = "text"
	}
	
	fieldContent = append(fieldContent,
		app.Input().
			Class("field-input").
			Type(inputType).
			Value(i.Value).
			Placeholder(i.Placeholder).
			Required(i.Required).
			Disabled(i.Disabled).
			OnChange(func(ctx app.Context, e app.Event) {
				if i.OnChange != nil {
					i.OnChange(e.Get("target").Get("value").String())
				}
			}),
	)
	
	if i.Error != "" {
		fieldContent = append(fieldContent,
			app.Div().Class("field-error-message").Text(i.Error),
		)
	}
	
	return app.Div().Class(className).Body(fieldContent...)
}

// StatusIndicator - 상태 표시 컴포넌트
type StatusIndicator struct {
	app.Compo
	Status  string // "online", "offline", "busy", "away"
	Size    string // "sm", "md", "lg"
	WithDot bool
}

func (s *StatusIndicator) Render() app.UI {
	classes := []string{"status-indicator"}
	classes = append(classes, fmt.Sprintf("status-%s", s.Status))
	
	if s.Size != "" {
		classes = append(classes, fmt.Sprintf("size-%s", s.Size))
	} else {
		classes = append(classes, "size-md")
	}
	
	if s.WithDot {
		classes = append(classes, "with-dot")
	}
	
	className := strings.Join(classes, " ")
	
	return app.Div().Class(className)
}

// SpreadsheetGrid is a spreadsheet-style grid component with Excel-like behavior.
// It provides a tabular interface with fixed column widths, cell selection,
// and support for different grid sizes and scrolling options.
//
// The component features:
//   - Fixed column widths to prevent layout shifting
//   - Excel-style column labels (A, B, C, etc.)
//   - Row numbering starting from 1
//   - Cell selection with visual feedback
//   - Text overflow handling with ellipsis
//   - Three predefined sizes (small, medium, large)
//   - Optional scrolling for mobile compatibility
//   - Special formatting for formulas, percentages, and currency
//
// Data Format:
// The Data map uses "row-col" keys where both row and col are zero-indexed.
// For example, "0-1" represents the cell at row 0, column 1 (cell B1 in Excel terms).
//
// Example usage:
//
//	data := map[string]string{
//	    "0-0": "Product", "0-1": "Price", "0-2": "Quantity",
//	    "1-0": "Apple",   "1-1": "$1.50", "1-2": "100",
//	    "2-0": "Orange",  "2-1": "$2.00", "2-2": "75",
//	}
//	
//	grid := &SpreadsheetGrid{
//	    Rows:       3,
//	    Cols:       3,
//	    Data:       data,
//	    Size:       "medium",
//	    Scrollable: true,
//	    MaxHeight:  "300px",
//	}
type SpreadsheetGrid struct {
	app.Compo
	
	// Rows specifies the number of data rows to display
	Rows int
	
	// Cols specifies the number of columns to display
	Cols int
	
	// Data contains the cell data using "row-col" format keys (e.g., "0-1" for row 0, column 1)
	Data map[string]string
	
	// Size determines the grid size: "small", "medium", or "large"
	// This affects column widths and font sizes
	Size string
	
	// Scrollable enables horizontal and vertical scrolling when content exceeds container
	Scrollable bool
	
	// MaxHeight sets the maximum height before vertical scrolling is enabled
	MaxHeight string
	
	// MaxWidth sets the maximum width before horizontal scrolling is enabled
	MaxWidth string
	
	// selectedCell tracks the currently selected cell in "row-col" format (internal use)
	selectedCell string
}

// OnMount initializes the SpreadsheetGrid component when it's mounted to the DOM.
// It ensures the Data map is properly initialized to prevent nil pointer errors.
func (g *SpreadsheetGrid) OnMount(ctx app.Context) {
	if g.Data == nil {
		g.Data = make(map[string]string)
	}
}

func (g *SpreadsheetGrid) Render() app.UI {
	classes := []string{"spreadsheet-grid"}
	classes = append(classes, fmt.Sprintf("grid-%s", g.Size))
	
	if g.Scrollable {
		classes = append(classes, "grid-scrollable")
	}
	
	className := strings.Join(classes, " ")
	
	// Build inline styles for scrollable grids
	var inlineStyle strings.Builder
	if g.Scrollable {
		if g.MaxHeight != "" {
			inlineStyle.WriteString(fmt.Sprintf("max-height: %s; ", g.MaxHeight))
		} else {
			inlineStyle.WriteString("max-height: 300px; ") // Default mobile height
		}
		if g.MaxWidth != "" {
			inlineStyle.WriteString(fmt.Sprintf("max-width: %s; ", g.MaxWidth))
		}
		inlineStyle.WriteString("overflow: auto; ")
	}
	
	var rows []app.UI
	
	// Header row with column labels
	headerCells := []app.UI{
		app.Div().Class("grid-cell header corner").Text(""),
	}
	for col := 0; col < g.Cols; col++ {
		headerCells = append(headerCells,
			app.Div().Class("grid-cell header").Text(g.getColumnLabel(col)),
		)
	}
	rows = append(rows, app.Div().Class("grid-row").Body(headerCells...))
	
	// Data rows
	for row := 0; row < g.Rows; row++ {
		var cells []app.UI
		
		// Row label
		cells = append(cells,
			app.Div().Class("grid-cell header").Text(fmt.Sprintf("%d", row+1)),
		)
		
		// Data cells
		for col := 0; col < g.Cols; col++ {
			cellKey := fmt.Sprintf("%d-%d", row, col)
			cellValue := g.Data[cellKey]
			
			cellClasses := []string{"grid-cell"}
			if g.selectedCell == cellKey {
				cellClasses = append(cellClasses, "selected")
			}
			
			// Special formatting for different types of content
			if strings.HasPrefix(cellValue, "=") {
				cellClasses = append(cellClasses, "formula")
			} else if strings.Contains(cellValue, "%") {
				cellClasses = append(cellClasses, "percentage")
			} else if strings.HasPrefix(cellValue, "$") {
				cellClasses = append(cellClasses, "currency")
			}
			
			cellClassName := strings.Join(cellClasses, " ")
			
			cells = append(cells,
				app.Div().
					Class(cellClassName).
					Text(cellValue).
					OnClick(func(ctx app.Context, e app.Event) {
						g.selectedCell = cellKey
						ctx.Update()
					}),
			)
		}
		
		rows = append(rows, app.Div().Class("grid-row").Body(cells...))
	}
	
	grid := app.Div().Class(className).Body(rows...)
	
	if g.Scrollable && inlineStyle.Len() > 0 {
		grid = grid.Style("style", inlineStyle.String())
	}
	
	return grid
}

// getColumnLabel converts a zero-based column index to Excel-style column labels.
// It generates labels like A, B, C, ..., Z, AA, AB, etc.
// For example: 0 -> "A", 25 -> "Z", 26 -> "AA", 27 -> "AB"
func (g *SpreadsheetGrid) getColumnLabel(col int) string {
	if col < 26 {
		return string(rune('A' + col))
	}
	return string(rune('A'+(col/26)-1)) + string(rune('A'+(col%26)))
}

// =============================================================================
// Database Table Grid Component (Embedding SpreadsheetGrid)
// =============================================================================

// ColumnMetadata represents metadata information for a database table column.
// It stores schema information such as data type, constraints, and descriptions
// that are typically found in database systems like PostgreSQL.
//
// Example usage:
//
//	column := ColumnMetadata{
//	    Name:        "user_id",
//	    Type:        "integer",
//	    Nullable:    false,
//	    Default:     "nextval('users_id_seq')",
//	    Description: "Primary key for users table",
//	}
type ColumnMetadata struct {
	// Name is the column name as it appears in the database
	Name string `json:"name"`
	
	// Type represents the SQL data type (e.g., "varchar", "integer", "decimal")
	Type string `json:"type"`
	
	// Nullable indicates whether the column allows NULL values
	Nullable bool `json:"nullable"`
	
	// Default specifies the default value for the column (optional)
	Default string `json:"default,omitempty"`
	
	// Description provides a human-readable description of the column (optional)
	Description string `json:"description,omitempty"`
	
	// Length specifies the maximum length for character types like varchar(50)
	Length int `json:"length,omitempty"`
	
	// Precision specifies the total number of digits for numeric types like decimal(10,2)
	Precision int `json:"precision,omitempty"`
	
	// Scale specifies the number of digits after decimal point for numeric types
	Scale int `json:"scale,omitempty"`
}

// DatabaseTableGrid is a PostgreSQL-style table grid component that embeds
// the functionality of SpreadsheetGrid while adding database-specific metadata display.
// It shows column information (name, type, nullable, etc.) above the actual data,
// similar to how database administration tools display table schemas.
//
// The component supports:
//   - Fixed column widths (Excel-like behavior)
//   - Text overflow with ellipsis
//   - Metadata rows showing column information
//   - Type-specific cell styling
//   - Horizontal and vertical scrolling
//   - Mobile-responsive design
//
// Example usage:
//
//	columns := []ColumnMetadata{
//	    {Name: "id", Type: "integer", Nullable: false, Description: "Primary key"},
//	    {Name: "name", Type: "varchar", Length: 100, Nullable: false},
//	}
//	
//	data := map[string]string{
//	    "0-0": "1", "0-1": "John Doe",
//	    "1-0": "2", "1-1": "Jane Smith",
//	}
//	
//	grid := &DatabaseTableGrid{
//	    TableName:       "users",
//	    Columns:         columns,
//	    Data:            data,
//	    ShowDescription: true,
//	    ShowConstraints: true,
//	    Scrollable:      true,
//	    MaxHeight:       "400px",
//	}
type DatabaseTableGrid struct {
	app.Compo
	
	// TableName is the name of the database table being displayed
	TableName string
	
	// Columns contains the metadata for each column in the table
	Columns []ColumnMetadata
	
	// Data contains the actual table data using "row-col" format keys (e.g., "0-1" for row 0, column 1)
	// This format is compatible with SpreadsheetGrid for easy data exchange
	Data map[string]string
	
	// MetadataRows specifies the number of metadata rows to display above the data
	// Default is 2 (field names and types), but increases with ShowDescription and ShowConstraints
	MetadataRows int
	
	// ShowDescription determines whether to display a description row for each column
	ShowDescription bool
	
	// ShowConstraints determines whether to display constraint information (defaults, etc.)
	ShowConstraints bool
	
	// Scrollable enables horizontal and vertical scrolling when content exceeds container size
	Scrollable bool
	
	// MaxHeight sets the maximum height before vertical scrolling is enabled
	MaxHeight string
	
	// MaxWidth sets the maximum width before horizontal scrolling is enabled  
	MaxWidth string
	
	// selectedCell tracks the currently selected cell in "row-col" format (internal use)
	selectedCell string
}

// OnMount initializes the DatabaseTableGrid component when it's mounted to the DOM.
// It ensures the Data map is initialized and calculates the appropriate number
// of metadata rows based on the ShowDescription and ShowConstraints flags.
func (d *DatabaseTableGrid) OnMount(ctx app.Context) {
	if d.Data == nil {
		d.Data = make(map[string]string)
	}
	
	// Calculate metadata rows based on enabled features
	if d.MetadataRows == 0 {
		d.MetadataRows = 2 // Base: field names and types
		if d.ShowDescription {
			d.MetadataRows++ // Add description row
		}
		if d.ShowConstraints {
			d.MetadataRows++ // Add constraints row
		}
	}
}

func (d *DatabaseTableGrid) Render() app.UI {
	// 기본 클래스
	classes := []string{"database-table-grid", "grid-scrollable"}
	className := strings.Join(classes, " ")
	
	// 스타일 빌더
	var inlineStyle strings.Builder
	if d.Scrollable {
		if d.MaxHeight != "" {
			inlineStyle.WriteString(fmt.Sprintf("max-height: %s; ", d.MaxHeight))
		} else {
			inlineStyle.WriteString("max-height: 400px; ")
		}
		if d.MaxWidth != "" {
			inlineStyle.WriteString(fmt.Sprintf("max-width: %s; ", d.MaxWidth))
		}
		inlineStyle.WriteString("overflow: auto; ")
	}
	
	var rows []app.UI
	
	// 테이블 제목
	if d.TableName != "" {
		rows = append(rows, 
			app.Div().Class("table-header").Body(
				app.H4().Class("table-name").Text(fmt.Sprintf("Table: %s", d.TableName)),
			),
		)
	}
	
	// 메타데이터 헤더 (항상 첫 번째 열은 "Field")
	metaHeaderCells := []app.UI{
		app.Div().Class("grid-cell meta-header corner").Text("Field"),
	}
	
	// 각 컬럼에 대한 헤더
	for _, col := range d.Columns {
		metaHeaderCells = append(metaHeaderCells,
			app.Div().Class("grid-cell meta-header").Text(col.Name),
		)
	}
	rows = append(rows, app.Div().Class("grid-row meta-row").Body(metaHeaderCells...))
	
	// 메타데이터 행들 생성
	rows = append(rows, d.renderMetadataRows()...)
	
	// 구분선
	separatorCells := make([]app.UI, len(d.Columns)+1)
	for i := range separatorCells {
		separatorCells[i] = app.Div().Class("grid-cell separator").Text("---")
	}
	rows = append(rows, app.Div().Class("grid-row separator-row").Body(separatorCells...))
	
	// 데이터 행들 (기존 SpreadsheetGrid 로직 활용)
	dataRows := d.getDataRowCount()
	for row := 0; row < dataRows; row++ {
		var cells []app.UI
		
		// 행 번호
		cells = append(cells,
			app.Div().Class("grid-cell row-number").Text(fmt.Sprintf("%d", row+1)),
		)
		
		// 데이터 셀들
		for col := 0; col < len(d.Columns); col++ {
			cellKey := fmt.Sprintf("%d-%d", row, col)
			cellValue := d.Data[cellKey]
			
			cellClasses := []string{"grid-cell", "data-cell"}
			if d.selectedCell == cellKey {
				cellClasses = append(cellClasses, "selected")
			}
			
			// 데이터 타입에 따른 스타일링
			if col < len(d.Columns) {
				switch d.Columns[col].Type {
				case "integer", "bigint", "decimal", "numeric":
					cellClasses = append(cellClasses, "cell-number")
				case "text", "varchar", "char":
					cellClasses = append(cellClasses, "cell-text")
				case "boolean":
					cellClasses = append(cellClasses, "cell-boolean")
				case "timestamp", "date", "time":
					cellClasses = append(cellClasses, "cell-datetime")
				}
			}
			
			cellClassName := strings.Join(cellClasses, " ")
			
			cells = append(cells,
				app.Div().
					Class(cellClassName).
					Text(cellValue).
					OnClick(func(ctx app.Context, e app.Event) {
						d.selectedCell = cellKey
						ctx.Update()
					}),
			)
		}
		
		rows = append(rows, app.Div().Class("grid-row data-row").Body(cells...))
	}
	
	grid := app.Div().Class(className).Body(rows...)
	
	if d.Scrollable && inlineStyle.Len() > 0 {
		grid = grid.Style("style", inlineStyle.String())
	}
	
	return grid
}

// renderMetadataRows generates the metadata rows that appear above the data.
// It creates rows for data types, nullable status, descriptions, and constraints
// based on the component configuration and column metadata.
func (d *DatabaseTableGrid) renderMetadataRows() []app.UI {
	var metaRows []app.UI
	
	// 타입 행
	typeCells := []app.UI{
		app.Div().Class("grid-cell meta-label").Text("Type"),
	}
	for _, col := range d.Columns {
		typeText := col.Type
		if col.Length > 0 {
			typeText = fmt.Sprintf("%s(%d)", col.Type, col.Length)
		} else if col.Precision > 0 && col.Scale > 0 {
			typeText = fmt.Sprintf("%s(%d,%d)", col.Type, col.Precision, col.Scale)
		}
		typeCells = append(typeCells,
			app.Div().Class("grid-cell meta-value type-cell").Text(typeText),
		)
	}
	metaRows = append(metaRows, app.Div().Class("grid-row meta-row").Body(typeCells...))
	
	// Nullable 행
	nullableCells := []app.UI{
		app.Div().Class("grid-cell meta-label").Text("Nullable"),
	}
	for _, col := range d.Columns {
		nullableText := "NO"
		if col.Nullable {
			nullableText = "YES"
		}
		nullableCells = append(nullableCells,
			app.Div().Class("grid-cell meta-value nullable-cell").Text(nullableText),
		)
	}
	metaRows = append(metaRows, app.Div().Class("grid-row meta-row").Body(nullableCells...))
	
	// 설명 행 (선택적)
	if d.ShowDescription {
		descCells := []app.UI{
			app.Div().Class("grid-cell meta-label").Text("Description"),
		}
		for _, col := range d.Columns {
			descText := col.Description
			if descText == "" {
				descText = "-"
			}
			descCells = append(descCells,
				app.Div().Class("grid-cell meta-value desc-cell").Text(descText),
			)
		}
		metaRows = append(metaRows, app.Div().Class("grid-row meta-row").Body(descCells...))
	}
	
	// 기본값 행 (선택적)
	if d.ShowConstraints {
		defaultCells := []app.UI{
			app.Div().Class("grid-cell meta-label").Text("Default"),
		}
		for _, col := range d.Columns {
			defaultText := col.Default
			if defaultText == "" {
				defaultText = "NULL"
			}
			defaultCells = append(defaultCells,
				app.Div().Class("grid-cell meta-value default-cell").Text(defaultText),
			)
		}
		metaRows = append(metaRows, app.Div().Class("grid-row meta-row").Body(defaultCells...))
	}
	
	return metaRows
}

// getDataRowCount calculates the number of data rows by parsing the Data map keys.
// It looks for the highest row number in "row-col" formatted keys and returns the total count.
func (d *DatabaseTableGrid) getDataRowCount() int {
	maxRow := -1
	for key := range d.Data {
		parts := strings.Split(key, "-")
		if len(parts) == 2 {
			var row int
			if _, err := fmt.Sscanf(parts[0], "%d", &row); err == nil && row > maxRow {
				maxRow = row
			}
		}
	}
	return maxRow + 1
}
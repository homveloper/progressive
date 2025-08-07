package grid

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// SpreadsheetGrid is a virtualized spreadsheet-style grid component with Excel-like behavior.
// It provides a tabular interface with infinite canvas support, virtualized rendering,
// and cell editing capabilities for optimal performance with large datasets.
//
// The component features:
//   - Virtualized rendering - only visible cells are rendered
//   - Infinite canvas support (default: 1M rows × 1K columns)
//   - Excel-style column labels (A, B, C, ..., AA, AB, etc.)
//   - Row numbering starting from 1
//   - Cell selection and inline editing with double-click
//   - Smooth scrolling performance with large datasets
//   - Sparse data storage - only stores cells with data
//   - Fixed column widths to prevent layout shifting
//   - Special formatting for formulas, percentages, and currency
//   - Optional scrolling for mobile compatibility
//
// Data Format:
// The Data map uses "row-col" keys where both row and col are zero-indexed.
// For example, "0-1" represents the cell at row 0, column 1 (cell B1 in Excel terms).
// Only cells with actual data are stored in the map (sparse storage).
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
//	    VirtualRows: 1000000, // 1M rows
//	    VirtualCols: 1000,    // 1K columns
//	    Data:        data,
//	    Size:        "medium",
//	    Scrollable:  true,
//	    MaxHeight:   "500px",
//	    CellWidth:   120,     // Fixed cell width for performance
//	    CellHeight:  32,      // Fixed cell height for performance
//	}
type SpreadsheetGrid struct {
	app.Compo

	// VirtualRows specifies the total virtual rows in the infinite canvas (default: 1000000)
	VirtualRows int

	// VirtualCols specifies the total virtual columns in the infinite canvas (default: 1000)
	VirtualCols int

	// Data contains the cell data using "row-col" format keys (e.g., "0-1" for row 0, column 1)
	// Only stores cells with actual data (sparse storage)
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

	// CellWidth sets the fixed width for each cell (default: 120px)
	CellWidth int

	// CellHeight sets the fixed height for each cell (default: 32px)
	CellHeight int

	// BufferSize sets how many extra rows/cols to render outside viewport (default: 5)
	BufferSize int

	// selectedCell tracks the currently selected cell in "row-col" format (internal use)
	selectedCell string

	// editingCell tracks the cell currently being edited (internal use)
	editingCell string

	// editingValue stores the temporary value while editing (internal use)
	editingValue string

	// viewport tracking (internal use)
	visibleStartRow int
	visibleEndRow   int
	visibleStartCol int
	visibleEndCol   int
	scrollTop       int
	scrollLeft      int

	// container dimensions (internal use)
	containerWidth  int
	containerHeight int
}

// OnMount initializes the SpreadsheetGrid component when it's mounted to the DOM.
// It ensures the Data map is properly initialized and sets up default values for virtualization.
func (g *SpreadsheetGrid) OnMount(ctx app.Context) {
	if g.Data == nil {
		g.Data = make(map[string]string)
	}

	// Set default values for infinite canvas
	if g.VirtualRows == 0 {
		g.VirtualRows = 1000000 // 1M rows like Excel
	}
	if g.VirtualCols == 0 {
		g.VirtualCols = 1000 // 1K columns for performance
	}
	if g.CellWidth == 0 {
		g.CellWidth = 120 // Default cell width
	}
	if g.CellHeight == 0 {
		g.CellHeight = 32 // Default cell height
	}
	if g.BufferSize == 0 {
		g.BufferSize = 5 // Buffer rows/cols for smooth scrolling
	}

	// Set default container dimensions
	if g.containerWidth == 0 {
		g.containerWidth = 800
	}
	if g.containerHeight == 0 {
		if g.MaxHeight != "" {
			// Try to parse MaxHeight and convert to pixels
			if strings.HasSuffix(g.MaxHeight, "px") {
				if h, err := strconv.Atoi(strings.TrimSuffix(g.MaxHeight, "px")); err == nil {
					g.containerHeight = h
				} else {
					g.containerHeight = 400 // fallback
				}
			} else {
				g.containerHeight = 400 // fallback
			}
		} else {
			g.containerHeight = 400 // default height
		}
	}

	// Calculate initial viewport
	g.calculateVisibleRange()
}

func (g *SpreadsheetGrid) Render() app.UI {
	classes := []string{"spreadsheet-grid", "virtualized-grid"}
	classes = append(classes, fmt.Sprintf("grid-%s", g.Size))

	if g.Scrollable {
		classes = append(classes, "grid-scrollable")
	}

	className := strings.Join(classes, " ")

	// Build inline styles for the grid container
	var containerStyle strings.Builder
	containerStyle.WriteString("position: relative; ")
	if g.Scrollable {
		if g.MaxHeight != "" {
			containerStyle.WriteString(fmt.Sprintf("height: %s; ", g.MaxHeight))
		} else {
			containerStyle.WriteString("height: 400px; ") // Default height
		}
		if g.MaxWidth != "" {
			containerStyle.WriteString(fmt.Sprintf("max-width: %s; ", g.MaxWidth))
		}
		containerStyle.WriteString("overflow: auto; ")
	}

	// Calculate total virtual dimensions
	totalWidth := g.VirtualCols * g.CellWidth + 60 // +60 for row header
	totalHeight := g.VirtualRows * g.CellHeight + 40 // +40 for column header

	// Create virtual scroll container
	virtualContainer := app.Div().
		Class("virtual-scroll-container").
		Style("width", fmt.Sprintf("%dpx", totalWidth)).
		Style("height", fmt.Sprintf("%dpx", totalHeight)).
		Style("position", "absolute").
		Style("pointer-events", "none")

	// Render only visible elements
	visibleElements := g.renderVisibleElements()

	// Main grid container with scroll handling
	gridContainer := app.Div().
		Class(className).
		Style("style", containerStyle.String()).
		OnScroll(g.handleScroll).
		Body(
			virtualContainer,
			visibleElements,
		)

	return gridContainer
}

// calculateVisibleRange determines which rows and columns are currently visible in the viewport
func (g *SpreadsheetGrid) calculateVisibleRange() {
	// Calculate visible rows
	headerHeight := 40
	visibleRowsCount := int(math.Ceil(float64(g.containerHeight-headerHeight) / float64(g.CellHeight)))
	g.visibleStartRow = int(math.Max(0, float64(g.scrollTop/g.CellHeight-g.BufferSize)))
	g.visibleEndRow = int(math.Min(float64(g.VirtualRows-1), float64(g.visibleStartRow+visibleRowsCount+2*g.BufferSize)))

	// Calculate visible columns
	headerWidth := 60
	visibleColsCount := int(math.Ceil(float64(g.containerWidth-headerWidth) / float64(g.CellWidth)))
	g.visibleStartCol = int(math.Max(0, float64(g.scrollLeft/g.CellWidth-g.BufferSize)))
	g.visibleEndCol = int(math.Min(float64(g.VirtualCols-1), float64(g.visibleStartCol+visibleColsCount+2*g.BufferSize)))
}

// handleScroll processes scroll events and updates the visible range
func (g *SpreadsheetGrid) handleScroll(ctx app.Context, e app.Event) {
	scrollTop := e.Get("target").Get("scrollTop").Int()
	scrollLeft := e.Get("target").Get("scrollLeft").Int()

	// Only update if scroll position changed significantly
	if math.Abs(float64(scrollTop-g.scrollTop)) > 10 || math.Abs(float64(scrollLeft-g.scrollLeft)) > 10 {
		g.scrollTop = scrollTop
		g.scrollLeft = scrollLeft
		g.calculateVisibleRange()
		ctx.Update()
	}
}

// renderVisibleElements creates the DOM elements for only the visible cells
func (g *SpreadsheetGrid) renderVisibleElements() app.UI {
	var elements []app.UI

	// Create positioned container for visible elements
	visibleContainer := app.Div().
		Class("visible-elements").
		Style("position", "absolute").
		Style("top", "0").
		Style("left", "0").
		Style("pointer-events", "auto")

	// Render column headers
	headerRow := g.renderColumnHeaders()
	elements = append(elements, headerRow)

	// Render visible data rows
	for row := g.visibleStartRow; row <= g.visibleEndRow; row++ {
		dataRow := g.renderDataRow(row)
		elements = append(elements, dataRow)
	}

	return visibleContainer.Body(elements...)
}

// renderColumnHeaders creates the column header row
func (g *SpreadsheetGrid) renderColumnHeaders() app.UI {
	var headerCells []app.UI

	// Corner cell
	cornerCell := app.Div().
		Class("grid-cell header corner").
		Style("position", "absolute").
		Style("top", "0").
		Style("left", "0").
		Style("width", "60px").
		Style("height", "40px").
		Style("z-index", "12").
		Text("")
	headerCells = append(headerCells, cornerCell)

	// Visible column headers
	for col := g.visibleStartCol; col <= g.visibleEndCol; col++ {
		left := 60 + col*g.CellWidth
		headerCell := app.Div().
			Class("grid-cell header").
			Style("position", "absolute").
			Style("top", "0").
			Style("left", fmt.Sprintf("%dpx", left)).
			Style("width", fmt.Sprintf("%dpx", g.CellWidth)).
			Style("height", "40px").
			Style("z-index", "11").
			Text(g.getColumnLabel(col))
		headerCells = append(headerCells, headerCell)
	}

	return app.Div().Body(headerCells...)
}

// renderDataRow creates a single data row with visible cells
func (g *SpreadsheetGrid) renderDataRow(row int) app.UI {
	var cells []app.UI
	top := 40 + row*g.CellHeight

	// Row header
	rowHeader := app.Div().
		Class("grid-cell header").
		Style("position", "absolute").
		Style("top", fmt.Sprintf("%dpx", top)).
		Style("left", "0").
		Style("width", "60px").
		Style("height", fmt.Sprintf("%dpx", g.CellHeight)).
		Style("z-index", "10").
		Text(fmt.Sprintf("%d", row+1))
	cells = append(cells, rowHeader)

	// Visible data cells in this row
	for col := g.visibleStartCol; col <= g.visibleEndCol; col++ {
		cell := g.renderCell(row, col, top)
		cells = append(cells, cell)
	}

	return app.Div().Body(cells...)
}

// renderCell creates a single cell with proper positioning and event handlers
func (g *SpreadsheetGrid) renderCell(row, col, top int) app.UI {
	cellKey := fmt.Sprintf("%d-%d", row, col)
	cellValue := g.Data[cellKey]
	left := 60 + col*g.CellWidth

	// Determine cell classes
	cellClasses := []string{"grid-cell"}
	if g.selectedCell == cellKey {
		cellClasses = append(cellClasses, "selected")
	}
	if g.editingCell == cellKey {
		cellClasses = append(cellClasses, "editing")
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

	// Create cell element
	cell := app.Div().
		Class(cellClassName).
		Style("position", "absolute").
		Style("top", fmt.Sprintf("%dpx", top)).
		Style("left", fmt.Sprintf("%dpx", left)).
		Style("width", fmt.Sprintf("%dpx", g.CellWidth)).
		Style("height", fmt.Sprintf("%dpx", g.CellHeight)).
		OnClick(func(ctx app.Context, e app.Event) {
			g.selectedCell = cellKey
			ctx.Update()
		}).
		OnDblClick(func(ctx app.Context, e app.Event) {
			g.startEditing(ctx, cellKey, cellValue)
		})

	// Render cell content based on editing state
	if g.editingCell == cellKey {
		// Render input field for editing
		input := app.Input().
			Type("text").
			Value(g.editingValue).
			Style("width", "100%").
			Style("height", "100%").
			Style("border", "none").
			Style("padding", "4px").
			Style("font-size", "inherit").
			Style("outline", "none").
			OnInput(func(ctx app.Context, e app.Event) {
				g.editingValue = e.Get("target").Get("value").String()
			}).
			OnKeyDown(func(ctx app.Context, e app.Event) {
				key := e.Get("key").String()
				switch key {
				case "Enter":
					g.finishEditing(ctx, true)
				case "Escape":
					g.finishEditing(ctx, false)
				}
			}).
			OnBlur(func(ctx app.Context, e app.Event) {
				g.finishEditing(ctx, true)
			})
		cell = cell.Body(input)
	} else {
		// Render cell text
		cell = cell.Text(cellValue)
	}

	return cell
}

// startEditing initiates cell editing mode
func (g *SpreadsheetGrid) startEditing(ctx app.Context, cellKey, currentValue string) {
	g.editingCell = cellKey
	g.editingValue = currentValue
	ctx.Update()
}

// finishEditing completes cell editing and optionally saves the value
func (g *SpreadsheetGrid) finishEditing(ctx app.Context, save bool) {
	if save && g.editingCell != "" {
		if g.editingValue != "" {
			g.Data[g.editingCell] = g.editingValue
		} else {
			// Remove cell from data if value is empty (sparse storage)
			delete(g.Data, g.editingCell)
		}
	}
	g.editingCell = ""
	g.editingValue = ""
	ctx.Update()
}

// getColumnLabel converts a zero-based column index to Excel-style column labels.
// It generates labels like A, B, C, ..., Z, AA, AB, etc.
// For example: 0 -> "A", 25 -> "Z", 26 -> "AA", 27 -> "AB"
func (g *SpreadsheetGrid) getColumnLabel(col int) string {
	if col < 26 {
		return string(rune('A' + col))
	}
	// Handle multi-character labels (AA, AB, etc.)
	result := ""
	for col >= 0 {
		result = string(rune('A'+(col%26))) + result
		col = col/26 - 1
		if col < 0 {
			break
		}
	}
	return result
}

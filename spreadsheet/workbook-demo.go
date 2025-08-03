package spreadsheet

import (
	"encoding/json"
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// WorkbookDemo demonstrates the Univer-inspired data model
type WorkbookDemo struct {
	app.Compo
	workbook     WorkbookData
	initialized  bool
	activeSheet  string
	selectedCell string
}

func (w *WorkbookDemo) OnMount(ctx app.Context) {
	// Create a demo workbook with sample data
	w.workbook = w.createDemoWorkbook()
	w.activeSheet = w.workbook.SheetOrder[0]
	w.initialized = true
}

func (w *WorkbookDemo) createDemoWorkbook() WorkbookData {
	// Create new workbook
	workbook := NewWorkbook("Financial Report 2024")

	// Get the default sheet and rename it
	defaultSheetID := workbook.SheetOrder[0]
	if sheet, ok := workbook.GetSheet(defaultSheetID); ok {
		sheet.Name = "Summary"
		sheet.TabColor = "#4285f4"
		workbook = workbook.UpdateSheet(defaultSheetID, sheet)
	}

	// Add some sample cells with different types
	// Headers
	workbook = workbook.UpdateCell(defaultSheetID, 0, 0, CellValue{
		Value: "Month",
		Type:  CellTypeString,
		Style: "header_style",
	})
	workbook = workbook.UpdateCell(defaultSheetID, 0, 1, CellValue{
		Value: "Revenue",
		Type:  CellTypeString,
		Style: "header_style",
	})
	workbook = workbook.UpdateCell(defaultSheetID, 0, 2, CellValue{
		Value: "Expenses",
		Type:  CellTypeString,
		Style: "header_style",
	})
	workbook = workbook.UpdateCell(defaultSheetID, 0, 3, CellValue{
		Value: "Profit",
		Type:  CellTypeString,
		Style: "header_style",
	})

	// Data rows
	months := []string{"January", "February", "March", "April"}
	revenues := []float64{15000, 18000, 16500, 20000}
	expenses := []float64{12000, 13500, 12800, 15000}

	for i, month := range months {
		row := i + 1

		// Month name
		workbook = workbook.UpdateCell(defaultSheetID, row, 0, CellValue{
			Value: month,
			Type:  CellTypeString,
		})

		// Revenue
		workbook = workbook.UpdateCell(defaultSheetID, row, 1, CellValue{
			Value: revenues[i],
			Type:  CellTypeNumber,
			Style: "currency_style",
		})

		// Expenses
		workbook = workbook.UpdateCell(defaultSheetID, row, 2, CellValue{
			Value: expenses[i],
			Type:  CellTypeNumber,
			Style: "currency_style",
		})

		// Profit (formula)
		workbook = workbook.UpdateCell(defaultSheetID, row, 3, CellValue{
			Value:   revenues[i] - expenses[i],
			Type:    CellTypeNumber,
			Formula: fmt.Sprintf("=B%d-C%d", row+1, row+1),
			Style:   "currency_style",
		})
	}

	// Add styles
	workbook.Styles["header_style"] = StyleData{
		Font: FontStyle{
			Bold:  true,
			Size:  12,
			Color: "#ffffff",
		},
		Background: BackgroundStyle{
			Color: "#4285f4",
		},
		Alignment: AlignmentStyle{
			Horizontal: "center",
		},
	}

	workbook.Styles["currency_style"] = StyleData{
		Format: NumberFormat{
			Pattern: "$#,##0.00",
			Type:    "currency",
		},
	}

	// Add a second sheet
	sheet2 := NewSheet("Details")
	sheet2.TabColor = "#34a853"
	sheet2.MergeData = append(sheet2.MergeData, MergeRange{
		StartRow: 0,
		StartCol: 0,
		EndRow:   0,
		EndCol:   3,
	})
	workbook = workbook.AddSheet(sheet2)

	// Add title to second sheet
	workbook = workbook.UpdateCell(sheet2.ID, 0, 0, CellValue{
		Value: "Quarterly Financial Details",
		Type:  CellTypeString,
		Style: "title_style",
	})

	workbook.Styles["title_style"] = StyleData{
		Font: FontStyle{
			Bold: true,
			Size: 16,
		},
		Alignment: AlignmentStyle{
			Horizontal: "center",
		},
	}

	return workbook
}

func (w *WorkbookDemo) Render() app.UI {
	// Check if workbook is initialized
	if !w.initialized {
		return app.Div().Class("workbook-demo").Body(
			app.H2().Text("Workbook Data Model Demo"),
			app.Div().Class("loading").Text("Loading workbook data..."),
		)
	}

	return app.Div().Class("workbook-demo").Body(
		app.H2().Text("Workbook Data Model Demo"),

		// Workbook info
		app.Div().Class("workbook-info").Body(
			app.H3().Text("Workbook Information"),
			app.Div().Class("info-grid").Body(
				app.Div().Body(
					app.Strong().Text("Name: "),
					app.Span().Text(w.workbook.Name),
				),
				app.Div().Body(
					app.Strong().Text("ID: "),
					app.Span().Text(w.workbook.ID),
				),
				app.Div().Body(
					app.Strong().Text("Sheets: "),
					app.Span().Text(fmt.Sprintf("%d", len(w.workbook.SheetOrder))),
				),
				app.Div().Body(
					app.Strong().Text("Version: "),
					app.Span().Text(w.workbook.AppVersion),
				),
			),
		),

		// Sheet tabs
		app.Div().Class("sheet-tabs-demo").Body(
			w.renderSheetTabs(),
		),

		// Active sheet content
		app.Div().Class("sheet-content").Body(
			w.renderActiveSheet(),
		),

		// Data structure view
		app.Div().Class("data-structure").Body(
			app.H3().Text("Data Structure (JSON)"),
			app.Pre().Class("json-view").Text(w.getJSONPreview()),
		),
	)
}

func (w *WorkbookDemo) renderSheetTabs() app.UI {
	if !w.initialized {
		return app.Div()
	}

	var tabs []app.UI

	for _, sheetID := range w.workbook.SheetOrder {
		sheet, _ := w.workbook.GetSheet(sheetID)
		tabClass := "demo-tab"
		if sheetID == w.activeSheet {
			tabClass += " active"
		}

		tabStyle := ""
		if sheet.TabColor != "" {
			tabStyle = fmt.Sprintf("border-bottom: 3px solid %s", sheet.TabColor)
		}

		tabs = append(tabs,
			app.Div().
				Class(tabClass).
				Style("style", tabStyle).
				Text(sheet.Name).
				OnClick(func(ctx app.Context, e app.Event) {
					w.activeSheet = sheetID
					ctx.Update()
				}),
		)
	}

	return app.Div().Class("tabs-container").Body(tabs...)
}

func (w *WorkbookDemo) renderActiveSheet() app.UI {
	if !w.initialized || w.activeSheet == "" {
		return app.Div()
	}

	sheet, ok := w.workbook.GetSheet(w.activeSheet)
	if !ok {
		return app.Div().Text("No sheet selected")
	}

	// Create a simplified grid view
	var rows []app.UI

	// Find the bounds of the data
	maxRow, maxCol := 0, 0
	for row := range sheet.CellData {
		if row > maxRow {
			maxRow = row
		}
		for col := range sheet.CellData[row] {
			if col > maxCol {
				maxCol = col
			}
		}
	}

	// Render grid
	for row := 0; row <= maxRow; row++ {
		var cells []app.UI

		for col := 0; col <= maxCol; col++ {
			cell, hasCell := sheet.GetCell(row, col)
			cellClass := "demo-cell"
			cellValue := ""
			cellStyle := ""

			if hasCell {
				cellValue = fmt.Sprintf("%v", cell.Value)

				// Apply style if exists
				if styleID, ok := cell.Style.(string); ok {
					if style, exists := w.workbook.Styles[styleID]; exists {
						if style.Background.Color != "" {
							cellStyle += fmt.Sprintf("background-color: %s; ", style.Background.Color)
						}
						if style.Font.Color != "" {
							cellStyle += fmt.Sprintf("color: %s; ", style.Font.Color)
						}
						if style.Font.Bold {
							cellStyle += "font-weight: bold; "
						}
						if style.Alignment.Horizontal != "" {
							cellStyle += fmt.Sprintf("text-align: %s; ", style.Alignment.Horizontal)
						}
					}
				}

				// Show formula indicator
				if cell.Formula != "" {
					cellClass += " has-formula"
					cellValue = cell.Formula
				}
			}

			cells = append(cells,
				app.Div().
					Class(cellClass).
					Style("style", cellStyle).
					Text(cellValue).
					Title(cellValue),
			)
		}

		rows = append(rows, app.Div().Class("demo-row").Body(cells...))
	}

	return app.Div().Body(
		app.H4().Text(fmt.Sprintf("Sheet: %s", sheet.Name)),
		app.Div().Class("demo-grid").Body(rows...),
	)
}

func (w *WorkbookDemo) getJSONPreview() string {
	if !w.initialized {
		return "{}"
	}

	// Create a simplified preview of the data structure
	preview := struct {
		Workbook struct {
			ID         string   `json:"id"`
			Name       string   `json:"name"`
			SheetOrder []string `json:"sheetOrder"`
		} `json:"workbook"`
		ActiveSheet struct {
			Name     string                 `json:"name"`
			CellData map[string]interface{} `json:"cellData"`
		} `json:"activeSheet"`
	}{}

	preview.Workbook.ID = w.workbook.ID
	preview.Workbook.Name = w.workbook.Name
	preview.Workbook.SheetOrder = w.workbook.SheetOrder

	if sheet, ok := w.workbook.GetSheet(w.activeSheet); ok {
		preview.ActiveSheet.Name = sheet.Name
		preview.ActiveSheet.CellData = make(map[string]interface{})

		// Show a few cells as example
		for row := 0; row <= 2; row++ {
			for col := 0; col <= 2; col++ {
				if cell, exists := sheet.GetCell(row, col); exists {
					key := fmt.Sprintf("cell[%d,%d]", row, col)
					preview.ActiveSheet.CellData[key] = map[string]interface{}{
						"value":   cell.Value,
						"type":    cell.Type,
						"formula": cell.Formula,
					}
				}
			}
		}
	}

	data, _ := json.MarshalIndent(preview, "", "  ")
	return string(data)
}

package main

import (
	"log"
	"net/http"
	"progressive/spreadsheet"
	"progressive/tetris"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Spreadsheet represents the main spreadsheet component
type Spreadsheet struct {
	app.Compo

	rows    int
	cols    int
	cells   map[string]string // key: "row-col", value: cell content
	editing string            // currently editing cell key
}

func (s *Spreadsheet) OnMount(ctx app.Context) {
	// Initialize with default grid size
	s.rows = 20
	s.cols = 10
	s.cells = make(map[string]string)
}

func (s *Spreadsheet) Render() app.UI {
	return app.Div().Class("spreadsheet-container").Body(
		app.H1().Text("Progressive Spreadsheet"),
		app.Div().Class("spreadsheet-grid").Body(
			s.renderGrid(),
		),
	)
}

func (s *Spreadsheet) renderGrid() app.UI {
	var rows []app.UI

	// Header row with column labels
	headerCells := []app.UI{
		app.Div().Class("cell header").Text(""), // top-left corner
	}
	for col := 0; col < s.cols; col++ {
		headerCells = append(headerCells,
			app.Div().Class("cell header").Text(s.getColumnLabel(col)),
		)
	}
	rows = append(rows, app.Div().Class("row").Body(headerCells...))

	// Data rows
	for row := 0; row < s.rows; row++ {
		var cells []app.UI

		// Row label
		cells = append(cells,
			app.Div().Class("cell header").Text(string(rune('1'+row))),
		)

		// Data cells
		for col := 0; col < s.cols; col++ {
			cellKey := s.getCellKey(row, col)
			cellValue := s.cells[cellKey]

			cells = append(cells, s.renderCell(row, col, cellKey, cellValue))
		}

		rows = append(rows, app.Div().Class("row").Body(cells...))
	}

	return app.Div().Body(rows...)
}

func (s *Spreadsheet) renderCell(row, col int, cellKey, value string) app.UI {
	isEditing := s.editing == cellKey

	if isEditing {
		return app.Input().
			Class("cell editing").
			Type("text").
			Value(value).
			AutoFocus(true).
			OnChange(s.ValueTo(&value)).
			OnBlur(func(ctx app.Context, e app.Event) {
				s.cells[cellKey] = value
				s.editing = ""
				ctx.Update()
			}).
			OnKeyDown(func(ctx app.Context, e app.Event) {
				if e.Get("key").String() == "Enter" {
					s.cells[cellKey] = value
					s.editing = ""
					ctx.Update()
				}
			})
	}

	return app.Div().
		Class("cell").
		Text(value).
		OnClick(func(ctx app.Context, e app.Event) {
			s.editing = cellKey
			ctx.Update()
		})
}

func (s *Spreadsheet) getCellKey(row, col int) string {
	return app.FormatString("%d-%d", row, col)
}

func (s *Spreadsheet) getColumnLabel(col int) string {
	if col < 26 {
		return string(rune('A' + col))
	}
	return string(rune('A'+(col/26)-1)) + string(rune('A'+(col%26)))
}

func main() {
	// Routes
	app.Route("/", func() app.Composer { return &spreadsheet.ComponentShowcasePage{} })
	app.Route("/simple", func() app.Composer { return &Spreadsheet{} })
	app.Route("/showcase", func() app.Composer { return &spreadsheet.ComponentShowcasePage{} })
	app.Route("/tetris", func() app.Composer { return &tetris.TetrisPage{} })

	// Start the app when running in the browser
	app.RunWhenOnBrowser()

	// HTTP handler for server-side
	http.Handle("/", &app.Handler{
		Name:        "Progressive Spreadsheet",
		Description: "A web-based spreadsheet built with Go and WebAssembly",
		Styles: []string{
			"/web/app.css",
			"/web/showcase.css",
			"/web/tetris.css",
		},
	})

	log.Println("Progressive Spreadsheet starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

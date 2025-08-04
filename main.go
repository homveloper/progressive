package main

import (
	"log"
	"net/http"
	"progressive/spreadsheet"
	"progressive/tetris"
	"progressive/tower"
	"progressive/qube"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	// Routes
	app.Route("/", func() app.Composer { return &spreadsheet.ComponentShowcasePage{} })
	app.Route("/showcase", func() app.Composer { return &spreadsheet.ComponentShowcasePage{} })
	app.Route("/tetris", func() app.Composer { return &tetris.TetrisPage{} })
	app.Route("/tower", func() app.Composer { return &tower.TowerDefensePage{} })
	app.Route("/qube", func() app.Composer { return &qube.QubePage{} })

	// Start the app when running in the browser
	app.RunWhenOnBrowser()

	// HTTP handler for server-side
	http.Handle("/", &app.Handler{
		Name:        "Progressive Games Collection",
		Description: "A collection of games built with Go and WebAssembly",
		Styles: []string{
			"/web/app.css",
			"/web/showcase.css",
			"/web/tetris.css",
			"/web/tower.css",
			"/web/qube.css",
		},
		Scripts: []string{
			"https://cdnjs.cloudflare.com/ajax/libs/three.js/r128/three.min.js",
			"/web/qube.js",
		},
	})

	log.Println("Progressive Games Collection starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

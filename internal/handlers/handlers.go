package handlers

import (
	"net/http"

	"progressive/internal/pages"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	component := pages.HomePage()
	component.Render(r.Context(), w)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	component := pages.Dashboard()
	component.Render(r.Context(), w)
}
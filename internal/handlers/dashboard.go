package handlers

import (
	"Praetor/internal/templates"
	"net/http"
)

type DashboardHandler struct{}

func (h *DashboardHandler) Page(w http.ResponseWriter, r *http.Request) {
	templates.Layout("Preator | Dashboard", templates.Dashboard()).Render(r.Context(), w)
}

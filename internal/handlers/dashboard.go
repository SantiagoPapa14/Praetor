package handlers

import (
	"Praetor/internal/app"
	"Praetor/internal/templates"
	"log"
	"net/http"
)

type DashboardHandler struct {
	App *app.App
}

func (h *DashboardHandler) Page(w http.ResponseWriter, r *http.Request) {
	templates.Layout("Preator | Dashboard", templates.Dashboard()).Render(r.Context(), w)
}

func (h *DashboardHandler) DockerTab(w http.ResponseWriter, r *http.Request) {
	containers, err := h.App.Repos.Docker.GetContainers()
	if err != nil {
		log.Println(err)
	}
	templates.DockerTab(containers).Render(r.Context(), w)
}

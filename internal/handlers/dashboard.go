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

func (h *DashboardHandler) DockerStart(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.App.Repos.Docker.StartContainer(id)
	if err != nil {
		log.Println(err)
	}

	cont, err := h.App.Repos.Docker.GetContainer(id)
	if err != nil {
		log.Println(err)
	}

	templates.DockerContainer(cont).Render(r.Context(), w)
}

func (h *DashboardHandler) DockerStop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.App.Repos.Docker.StopContainer(id)
	if err != nil {
		log.Println(err)
	}

	cont, err := h.App.Repos.Docker.GetContainer(id)
	if err != nil {
		log.Println(err)
	}

	templates.DockerContainer(cont).Render(r.Context(), w)
}

func (h *DashboardHandler) DockerRestart(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.App.Repos.Docker.RestartContainer(id)
	if err != nil {
		log.Println(err)
	}

	cont, err := h.App.Repos.Docker.GetContainer(id)
	if err != nil {
		log.Println(err)
	}

	templates.DockerContainer(cont).Render(r.Context(), w)
}

func (h *DashboardHandler) DockerLogs(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logs, err := h.App.Repos.Docker.GetContainerLogs(id, "500")
	if err != nil {
		log.Println(err)
	}

	templates.DockerLogs(logs).Render(r.Context(), w)
}

package handlers

import (
	"Praetor/internal/models"
	"Praetor/internal/repositories"
	"Praetor/internal/templates"
	"net/http"
	"strconv"
)

type PhraseHandler struct {
	Repository *repositories.PhraseRepository
}

func (h *PhraseHandler) Page(w http.ResponseWriter, r *http.Request) {
	phrases, err := h.Repository.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	templates.Layout("Phrases", templates.PhrasesPage(phrases)).Render(r.Context(), w)
}

func (h *PhraseHandler) List(w http.ResponseWriter, r *http.Request) {
	phrases, err := h.Repository.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	templates.PhraseList(phrases).Render(r.Context(), w)
}

func (h *PhraseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	phrase, err := h.Repository.GetByID(id)

	if phrase == nil {
		http.Error(w, "Phrase not found", 404)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = h.Repository.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/phrases", http.StatusSeeOther)
}

func (h *PhraseHandler) Add(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	phrase, err := h.Repository.Create(text)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	templates.Phrase(models.Phrase(*phrase)).Render(r.Context(), w)
}

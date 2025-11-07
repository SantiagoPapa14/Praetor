package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"Praetor/internal/models"
	"Praetor/internal/templates"
)

type PhraseHandler struct {
	DB *sql.DB
}

func (h *PhraseHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, text FROM phrases ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var phrases []models.Phrase
	for rows.Next() {
		var text string
		var id int
		if err := rows.Scan(&id, &text); err == nil {
			phrases = append(phrases, models.Phrase{ID: id, Text: text})
		}
	}

	templates.Layout("Phrases", templates.Phrases(phrases)).Render(r.Context(), w)
}

func (h *PhraseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Println("Deleting phrase", id)
	_, err := h.DB.Exec("DELETE FROM phrases WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Return updated list fragment (for HTMX)
	rows, err := h.DB.Query("SELECT id, text FROM phrases ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var phrases []models.Phrase
	for rows.Next() {
		var text string
		var id int
		if err := rows.Scan(&id, &text); err == nil {
			phrases = append(phrases, models.Phrase{ID: id, Text: text})
		}
	}

	templates.PhraseList(phrases).Render(r.Context(), w)
}

func (h *PhraseHandler) Add(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	if text != "" {
		_, err := h.DB.Exec("INSERT INTO phrases (text) VALUES (?)", text)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	// Return updated list fragment (for HTMX)
	rows, err := h.DB.Query("SELECT id, text FROM phrases ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var phrases []models.Phrase
	for rows.Next() {
		var text string
		var id int
		if err := rows.Scan(&id, &text); err == nil {
			phrases = append(phrases, models.Phrase{ID: id, Text: text})
		}
	}

	templates.PhraseList(phrases).Render(r.Context(), w)
}

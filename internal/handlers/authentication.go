package handlers

import (
	"Praetor/internal/auth"
	"Praetor/internal/repositories"
	"Praetor/internal/templates"
	"net/http"
)

type AuthenticationHandler struct {
	Repository *repositories.SessionRepository
}

func (h *AuthenticationHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	templates.Layout("Macaco Login", templates.Login()).Render(r.Context(), w)
}

func (h *AuthenticationHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Redirect(w, r, "/authenticate", http.StatusSeeOther)
		return
	}

	if err := auth.Login(*h.Repository, w, r, email, password); err != nil {
		http.Redirect(w, r, "/authenticate", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

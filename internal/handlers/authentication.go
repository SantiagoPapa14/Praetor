package handlers

import (
	"Praetor/internal/auth"
	"Praetor/internal/models"
	"Praetor/internal/repositories"
	"Praetor/internal/templates"
	"log"
	"net/http"
	"time"
)

type AuthenticationHandler struct {
	SessionRepository *repositories.SessionRepository
	UserRepository    *repositories.UserRepository
}

func (h *AuthenticationHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	templates.Layout("Praetor | Login", templates.Login()).Render(r.Context(), w)
}

func (h *AuthenticationHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		templates.AuthResponseMessage("Email and Password required").Render(r.Context(), w)
		return
	}

	user, err := h.UserRepository.GetByEmail(email)
	if err != nil {
		log.Println(err)
		templates.AuthResponseMessage("Invalid email or password").Render(r.Context(), w)
		return
	}

	if password != user.Password {
		log.Println(err)
		templates.AuthResponseMessage("Invalid email or password").Render(r.Context(), w)
		return
	}

	err = auth.CreateSession(*h.SessionRepository, w, user.ID, 24*time.Hour)
	if err != nil {
		log.Println(err)
		templates.AuthResponseMessage("Error logging you in").Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthenticationHandler) Register(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm_password := r.FormValue("confirm_password")

	if name == "" || email == "" || password == "" || confirm_password == "" {
		templates.AuthResponseMessage("All fields are required").Render(r.Context(), w)
		return
	}

	if password != confirm_password {
		templates.AuthResponseMessage("Passwords do not match").Render(r.Context(), w)
		return
	}

	var userToCreate models.User = models.User{
		Email:    email,
		Password: password,
	}

	err := h.UserRepository.Create(&userToCreate)
	if err != nil {
		templates.AuthResponseMessage("Error registering user").Render(r.Context(), w)
		return
	}

	templates.AuthResponseMessage("User registered successfully").Render(r.Context(), w)
}

func (h *AuthenticationHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		err = auth.DeleteSession(*h.SessionRepository, cookie.Value)
		if err != nil {
			log.Println(err)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // <-- delete
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // true in prod
	})

	w.Header().Set("HX-Redirect", "/authenticate")
	w.WriteHeader(http.StatusOK)
}

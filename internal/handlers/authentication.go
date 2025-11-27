package handlers

import (
	"Praetor/internal/app"
	"Praetor/internal/auth"
	"Praetor/internal/models"
	"Praetor/internal/templates"
	"log"
	"net/http"
	"time"
)

type AuthenticationHandler struct {
	App *app.App
}

func (h *AuthenticationHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	templates.Layout("Praetor | Login", templates.Login()).Render(r.Context(), w)
}

func (h *AuthenticationHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		templates.AuthResponseMessage("Email and Password required", "login-response-message").Render(r.Context(), w)
		return
	}

	user, err := h.App.Repos.User.GetByEmail(email)
	if err != nil {
		log.Println(err)
		templates.AuthResponseMessage("Invalid email or password", "login-response-message").Render(r.Context(), w)
		return
	}

	if password != user.Password {
		log.Println(err)
		templates.AuthResponseMessage("Invalid email or password", "login-response-message").Render(r.Context(), w)
		return
	}

	err = auth.CreateSession(*h.App.Repos.Session, w, user.ID, 24*time.Hour)
	if err != nil {
		log.Println(err)
		templates.AuthResponseMessage("Error logging you in", "login-response-message").Render(r.Context(), w)
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
		templates.AuthResponseMessage("All fields are required", "register-response-message").Render(r.Context(), w)
		return
	}

	if password != confirm_password {
		templates.AuthResponseMessage("Passwords do not match", "register-response-message").Render(r.Context(), w)
		return
	}

	var userToCreate models.User = models.User{
		Email:    email,
		Password: password,
	}

	err := h.App.Repos.User.Create(&userToCreate)
	if err != nil {
		templates.AuthResponseMessage("Error registering user", "register-response-message").Render(r.Context(), w)
		return
	}

	templates.AuthCustomResponseMessage("User registered successfully", "register-response-message", "text-center text-green-500 text-sm mt-2").Render(r.Context(), w)
}

func (h *AuthenticationHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		err = auth.DeleteSession(*h.App.Repos.Session, cookie.Value)
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

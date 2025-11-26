package middleware

import (
	"Praetor/internal/repositories"
	"context"
	"net/http"
	"time"
)

func AuthMiddleware(sessionRepo *repositories.SessionRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/authenticate", http.StatusSeeOther)
			return
		}

		session, err := sessionRepo.GetByToken(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/authenticate", http.StatusSeeOther)
			return
		}

		expiry, err := time.Parse(time.RFC3339, session.ExpiresAt)
		if err != nil || expiry.Before(time.Now().UTC()) {
			_ = sessionRepo.Delete(cookie.Value)
			http.Redirect(w, r, "/authenticate", http.StatusSeeOther)
			return
		}

		sessionRepo.UpdateLastSeen(cookie.Value)

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", session.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

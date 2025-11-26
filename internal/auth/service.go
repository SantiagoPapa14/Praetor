package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"Praetor/internal/models"
	"Praetor/internal/repositories"
)

func GenerateToken(size int) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func createSession(repository repositories.SessionRepository, w http.ResponseWriter, userID int, duration time.Duration) error {
	token, err := GenerateToken(32)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	expires := now.Add(duration)

	session := &models.Session{
		Token:     token,
		UserID:    userID,
		CreatedAt: now.Format(time.RFC3339),
		LastSeen:  now.Format(time.RFC3339),
		ExpiresAt: expires.Format(time.RFC3339),
	}

	err = repository.Create(session)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // SET TRUE IN PRODUCTION
		Expires:  expires,
	}
	http.SetCookie(w, cookie)
	return nil
}

func Login(repository *repositories.SessionRepository, username string, password string) error {
	return createSession(*repository, w, userID, 24*time.Hour)
}

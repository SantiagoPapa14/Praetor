package repositories

import (
	"Praetor/internal/models"
	"database/sql"
	"errors"
	"time"
)

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) Create(session *models.Session) error {
	_, err := r.DB.Exec("INSERT INTO sessions (token, user_id, created_at, last_seen_at, expires_at) VALUES (?, ?, ?, ?, ?)",
		session.Token, session.UserID, session.CreatedAt, session.LastSeen, session.ExpiresAt)
	return err
}

func (r *SessionRepository) Delete(token string) error {
	_, err := r.DB.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}

func (r *SessionRepository) GetByToken(token string) (*models.Session, error) {
	var session models.Session
	err := r.DB.QueryRow("SELECT token, user_id, created_at, last_seen_at, expires_at FROM sessions WHERE token = ?", token).
		Scan(&session.Token, &session.UserID, &session.CreatedAt, &session.LastSeen, &session.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("Session not found")
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) UpdateLastSeen(token string) error {
	lastSeen := time.Now().UTC().Format(time.RFC3339)
	_, err := r.DB.Exec("UPDATE sessions SET last_seen_at = ? WHERE token = ?", lastSeen, token)
	return err
}

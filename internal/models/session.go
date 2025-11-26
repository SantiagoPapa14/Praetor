package models

type Session struct {
	Token     string
	UserID    int
	CreatedAt string
	LastSeen  string
	ExpiresAt string
}

package repositories

import (
	"Praetor/internal/models"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	_, err := r.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		user.Name, user.Email, user.Password)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	row := r.DB.QueryRow("SELECT id, email, name, password FROM users WHERE email = ?", email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

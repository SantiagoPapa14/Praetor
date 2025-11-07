package repositories

import (
	"Praetor/internal/models"
	"database/sql"
)

type PhraseRepository struct {
	DB *sql.DB
}

func NewPhraseRepository(db *sql.DB) *PhraseRepository {
	return &PhraseRepository{DB: db}
}

// GetAll retrieves all phrases ordered by ID descending
func (r *PhraseRepository) GetAll() ([]models.Phrase, error) {
	rows, err := r.DB.Query("SELECT id, text FROM phrases ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phrases []models.Phrase
	for rows.Next() {
		var phrase models.Phrase
		if err := rows.Scan(&phrase.ID, &phrase.Text); err != nil {
			return nil, err
		}
		phrases = append(phrases, phrase)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return phrases, nil
}

// GetByID retrieves a single phrase by ID
func (r *PhraseRepository) GetByID(id int) (*models.Phrase, error) {
	var phrase models.Phrase
	err := r.DB.QueryRow("SELECT id, text FROM phrases WHERE id = ?", id).
		Scan(&phrase.ID, &phrase.Text)

	if err == sql.ErrNoRows {
		return nil, nil // or return a custom "not found" error
	}
	if err != nil {
		return nil, err
	}

	return &phrase, nil
}

// Create inserts a new phrase
func (r *PhraseRepository) Create(text string) (*models.Phrase, error) {
	result, err := r.DB.Exec("INSERT INTO phrases (text) VALUES (?)", text)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	createdPhrase := models.Phrase{
		ID:   int(id),
		Text: text,
	}

	return &createdPhrase, nil
}

// Update modifies an existing phrase
func (r *PhraseRepository) Update(id int, text string) error {
	_, err := r.DB.Exec("UPDATE phrases SET text = ? WHERE id = ?", text, id)
	return err
}

// Delete removes a phrase by ID
func (r *PhraseRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM phrases WHERE id = ?", id)
	return err
}

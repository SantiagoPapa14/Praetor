package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:data.db?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS phrases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL
	);`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return db, nil
}

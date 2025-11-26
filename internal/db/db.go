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
	);

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP)
	);

	CREATE TABLE IF NOT EXISTS sessions (
  		token TEXT PRIMARY KEY,
  		user_id INTEGER NOT NULL,
  		created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  		last_seen_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  		expires_at TEXT NOT NULL,

		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return db, nil
}

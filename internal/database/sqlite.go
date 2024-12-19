package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type DB struct {
	*sql.DB
}

type Session struct {
	ID           string
	SessionToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func NewSQLiteDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		session_token TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	return err
}

func (db *DB) SaveSession(session *Session) error {
	query := `
	INSERT OR REPLACE INTO sessions (id, session_token, expires_at)
	VALUES (?, ?, ?)`

	_, err := db.Exec(query, session.ID, session.SessionToken, session.ExpiresAt)
	return err
}

func (db *DB) GetLatestSession() (*Session, error) {
	query := `
	SELECT id, session_token, expires_at, created_at
	FROM sessions
	ORDER BY created_at DESC
	LIMIT 1`

	var session Session
	err := db.QueryRow(query).Scan(
		&session.ID,
		&session.SessionToken,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
} 
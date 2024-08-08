package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var connStr = "postgres://postgres:pass123@postgres:5432/event_streaming?sslmode=disable"

func InitDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return db, nil
}

func FetchLatestEvent(db *sql.DB) (string, error) {
	var event string
	query := `SELECT * FROM events ORDER BY id DESC LIMIT 1`

	row := db.QueryRow(query)
	err := row.Scan(&event)
	if err != nil {
		return "", err
	}

	return event, nil
}

func CreateEvent(db *sql.DB, event string) (int, error) {
	query := `INSERT INTO events (event) VALUES ($1) RETURNING event_id`

	var eventID int
	err := db.QueryRow(query, event).Scan(&eventID)
	if err != nil {
			return 0, err
	}

	return eventID, nil
}
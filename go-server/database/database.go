package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var connStr = "postgres://postgres:pass123@postgres:5432/event_streaming?sslmode=disable"

type Event struct {
	EventID           int
	GuestSatisfaction bool
	StressMarks       int
	Timestamp         string
}

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

func FetchLatestEvent(db *sql.DB) (int, error) {
	var event_id int
	query := `SELECT event_id FROM events ORDER BY id DESC LIMIT 1`

	row := db.QueryRow(query)
	err := row.Scan(&event_id)
	if err != nil {
		return 0, err
	}

	return event_id, nil
}

func CreateEvent(db *sql.DB, event Event) (int, error) {
	query := `INSERT INTO events (event) VALUES ($1) RETURNING event_id`

	var eventID int
	err := db.QueryRow(query, event).Scan(&eventID)
	if err != nil {
			return 0, err
	}

	return eventID, nil
}
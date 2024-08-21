package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	var timestamp time.Time

	query := `SELECT event_id, timestamp FROM events ORDER BY event_id DESC LIMIT 1`
	row := db.QueryRow(query)
	err := row.Scan(&event_id, &timestamp)
	if err == sql.ErrNoRows {
			log.Println("No events found in the database.")
			return 0, nil
	} else if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			return 0, err
	}

	// Check if the timestamp is greater than 6 minutes ago
	if timestamp.Before(time.Now().Add(-6 * time.Minute)) {
		log.Println("Latest event is older than 6 minutes.")
		return 0, nil
	}

	return event_id, nil
}

func CreateEvent(db *sql.DB, event Event) (int, error) {
	query := `INSERT INTO events (guest_satisfaction, stress_marks) VALUES ($1, $2) RETURNING event_id`

	var eventID int
	err := db.QueryRow(query, event.GuestSatisfaction, event.StressMarks).Scan(&eventID)
	if err != nil {
			return 0, err
	}

	return eventID, nil
}
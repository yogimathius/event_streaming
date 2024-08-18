package db

import (
	"database/sql"
	"fmt"
	"go-consumer/message"
	"log"

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

func FetchWorkerData(db *sql.DB, teamID int) (string, error) {
	var routineType string
	query := `SELECT workers FROM event_streaming WHERE team_id = $1`

	row := db.QueryRow(query, teamID)
	err := row.Scan(&routineType)
	if err != nil {
		return "", err
	}

	return routineType, nil
}

func AddEventMessage(db *sql.DB, eventID int, message message.Message) error {
	query := `INSERT INTO event_messages (event_id, event_type, priority, description, status) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, eventID, message.EventType, message.Priority, message.Description, message.Status)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

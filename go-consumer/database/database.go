package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var connStr = "postgres://postgres:password@localhost:5432/event_streaming"

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

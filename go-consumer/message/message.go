package message

import "time"

type Message struct {
	Id          string    `json:"id"`
	EventId 	  int       `json:"event_id"`
	EventType   string    `json:"event_type"`
	EventTime   time.Time `json:"event_time"`
	Priority    string    `json:"priority"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
}

package ws

import (
	"encoding/json"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeTaskCreated EventType = "task_created"
)

type Event struct {
	EventType EventType       `json:"type"`
	Data      json.RawMessage `json:"data"`
}

func (e Event) AsEventTaskCreated() (EventTaskCreated, error) {
	var data EventTaskCreated
	err := json.Unmarshal(e.Data, &data)
	return data, err
}

type EventTaskCreated struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedBy uint      `json:"created_by"`
}

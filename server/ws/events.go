package ws

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeTaskCreated      EventType = "task_created"
	EventTypeTaskStoreFailure EventType = "task_store_error"
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

func FromEventStoreFailure(data EventTaskStoreFailure) (Event, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return Event{}, fmt.Errorf("failed to marshal event: %w", err)
	}

	return Event{
		EventType: EventTypeTaskStoreFailure,
		Data:      body,
	}, nil
}

type EventTaskCreated struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedBy uint      `json:"created_by"`
}

type EventTaskStoreFailure struct {
	Error string `json:"error"`
}

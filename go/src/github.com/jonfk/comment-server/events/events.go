package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

const (
	AccountCreatedTypeName       = "AccountCreatedEvent"
	AccountDeletedTypeName       = "AccountDeleted"
	CommentThreadCreatedTypeName = "CommentThreadCreated"
	CommentCreatedTypeName       = "CommentCreated"
	CommentDeletedTypeName       = "CommentDeleted"
)

type EventPayload interface {
	EventType() string
}

type AccountCreated struct {
	AccountId uuid.UUID `json:"accountId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
}

type AccountDeleted struct {
	AccountId uuid.UUID `json:"accountId"`
}

type CommentThreadCreated struct {
	CommentThreadId uuid.UUID `json:"commentThreadId"`
	PageUrl         string    `json:"pageUrl"`
	Title           string    `json:"title"`
}

type CommentCreated struct {
	CommentId       uuid.UUID  `json:"commentId"`
	Data            string     `json:"data"`
	ParentId        *uuid.UUID `json:"parentId"`
	CommentThreadId uuid.UUID  `json:"commentThreadId"`
	AccountId       uuid.UUID  `json:"accountId"`
}

type CommentDeleted struct {
	CommentId uuid.UUID `json:"commentId"`
}

func (e AccountCreated) EventType() string       { return AccountCreatedTypeName }
func (e AccountDeleted) EventType() string       { return AccountDeletedTypeName }
func (e CommentThreadCreated) EventType() string { return CommentThreadCreatedTypeName }
func (e CommentCreated) EventType() string       { return CommentCreatedTypeName }
func (e CommentDeleted) EventType() string       { return CommentDeletedTypeName }

type Event struct {
	EventType string       `json:"eventType"`
	Timestamp time.Time    `json:"timestamp"`
	Payload   EventPayload `json:"payload"`
}

type EventJSON struct {
	EventType string          `json:"eventType"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

func NewEvent(timestamp time.Time, payload EventPayload) Event {
	return Event{EventType: payload.EventType(),
		Timestamp: timestamp,
		Payload:   payload}
}

func MarshalJSON(event Event) ([]byte, error) {
	return json.Marshal(event)
}

func UnmarshalJSON(input []byte) (Event, error) {
	rawEvent := EventJSON{}
	err := json.Unmarshal(input, &rawEvent)
	if err != nil {
		return Event{}, err
	}

	switch rawEvent.EventType {
	case AccountCreatedTypeName:
		eventPayload := AccountCreated{}
		err = json.Unmarshal(rawEvent.Payload, &eventPayload)
		if err != nil {
			return Event{}, err
		}
		return Event{EventType: eventPayload.EventType(),
			Timestamp: rawEvent.Timestamp,
			Payload:   eventPayload}, nil
	case AccountDeletedTypeName:
		eventPayload := AccountDeleted{}
		err = json.Unmarshal(rawEvent.Payload, &eventPayload)
		if err != nil {
			return Event{}, err
		}
		return Event{EventType: eventPayload.EventType(),
			Timestamp: rawEvent.Timestamp,
			Payload:   eventPayload}, nil
	case CommentThreadCreatedTypeName:
		eventPayload := CommentThreadCreated{}
		err = json.Unmarshal(rawEvent.Payload, &eventPayload)
		if err != nil {
			return Event{}, err
		}
		return Event{EventType: eventPayload.EventType(),
			Timestamp: rawEvent.Timestamp,
			Payload:   eventPayload}, nil
	case CommentCreatedTypeName:
		eventPayload := CommentCreated{}
		err = json.Unmarshal(rawEvent.Payload, &eventPayload)
		if err != nil {
			return Event{}, err
		}
		return Event{EventType: eventPayload.EventType(),
			Timestamp: rawEvent.Timestamp,
			Payload:   eventPayload}, nil
	case CommentDeletedTypeName:
		eventPayload := CommentDeleted{}
		err = json.Unmarshal(rawEvent.Payload, &eventPayload)
		if err != nil {
			return Event{}, err
		}
		return Event{EventType: eventPayload.EventType(),
			Timestamp: rawEvent.Timestamp,
			Payload:   eventPayload}, nil
	default:
		return Event{}, fmt.Errorf("unknown event type %s", rawEvent.EventType)
	}
}

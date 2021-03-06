package events

import (
	"reflect"
	"testing"

	"github.com/satori/go.uuid"
)

func TestMarshalJSON(t *testing.T) {
	parentId := uuid.NewV4()
	eventPayloads := []EventPayload{
		AccountCreated{AccountId: uuid.NewV4(),
			Username:       "username",
			Email:          "email@example.com",
			HashedPassword: []byte("hashed_password"),
			HashSalt:       []byte("scrypt salt")},
		AccountDeleted{AccountId: uuid.NewV4()},
		AccountLoggedIn{AccountId: uuid.NewV4(), JWT: "invalid_token"},
		CommentThreadCreated{CommentThreadId: uuid.NewV4(),
			PageUrl: "pageurl.com",
			Title:   "title"},
		CommentCreated{CommentId: uuid.NewV4(),
			Data:            "this is a comment",
			ParentId:        &parentId,
			CommentThreadId: uuid.NewV4(),
			AccountId:       uuid.NewV4()},
		CommentDeleted{CommentId: uuid.NewV4()},
	}

	expectedEvents := []Event{}

	for _, payload := range eventPayloads {
		event := NewEventNow(payload)
		expectedEvents = append(expectedEvents, event)
	}

	encodedEvents := [][]byte{}
	decodedEvents := []Event{}

	for _, event := range expectedEvents {
		encodedEvent, err := MarshalJSON(event)
		if err != nil {
			t.Fatalf("MarshalJSON failed : %v", err)
		}
		//fmt.Println(string(encodedEvent))
		encodedEvents = append(encodedEvents, encodedEvent)
	}

	for _, encodedEvent := range encodedEvents {
		decodedEvent, err := UnmarshalJSON(encodedEvent)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed : %v", err)
		}
		decodedEvents = append(decodedEvents, decodedEvent)
	}

	for i, decodedEvent := range decodedEvents {
		if !reflect.DeepEqual(expectedEvents[i], decodedEvent) {
			t.Fatalf("expectedEvent != decodedEvent\n(expectedEvent) %v != (decodedEvent) %v", expectedEvents[i], decodedEvent)
		}
	}
}

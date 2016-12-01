package accounts

import (
	"fmt"
	"github.com/jonfk/comment-server/events"
)

type EventHandler struct {
	AccountsService *Accounts
}

func (e *EventHandler) HandleEvent(event events.Event) error {

	switch eventPayload := event.Payload.(type) {
	case events.AccountCreated:
	case events.AccountDeleted:
	case events.CommentThreadCreated:
	case events.CommentCreated:
	case events.CommentDeleted:
	default:
		fmt.Println(eventPayload)
	}
	return nil
}

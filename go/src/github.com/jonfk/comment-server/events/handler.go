package events

import (
	log "github.com/Sirupsen/logrus"
)

type EventHandler interface {
	HandleEvent(Event)
}

// A Default implementation of EventHandler that simply logs the event
// being handled
type LogEventHandler struct{}

func (handler LogEventHandler) HandleEvent(event Event) {
	log.WithFields(log.Fields{
		"context": "LogEventHandler",
		"event":   event,
	}).Info("Event Handled")
}

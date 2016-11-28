package commands

import (
	"github.com/jonfk/comment-server/events"
)

// a CommandHandler handles a Command.
// It can interpret the command and return an event and a nil error
// or an error.
//
// Careful about how commands are logged. Commands can contain
// sensitive information that shouldn't be store such as unhashedPasswords.
type CommandHandler interface {
	HandleCommand(Command) (events.Event, error)
}

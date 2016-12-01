package accounts

import (
	"fmt"
	"time"

	"github.com/jonfk/comment-server/commands"
	"github.com/jonfk/comment-server/events"

	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

type CommandHandler struct {
	EventHandler    events.EventHandler
	AccountsService *Accounts
}

func (c *CommandHandler) HandleCommand(command commands.Command) (events.Event, error) {
	switch commandPayload := command.Payload.(type) {
	case commands.CreateAccount:
		salt, err := GenerateSalt()
		if err != nil {
			// Fix error to be friendly
			return events.Event{}, err
		}
		hashedPassword, err := HashPassword(commandPayload.UnhashedPassword, salt)
		if err != nil {
			// Fix error to be friendly
			return events.Event{}, err
		}

		eventPayload := events.AccountCreated{
			AccountId:      uuid.NewV4(),
			Username:       commandPayload.Username,
			Email:          commandPayload.Email,
			HashedPassword: hashedPassword,
			HashSalt:       salt,
		}
		event := events.NewEvent(time.Now().UTC().Round(time.Second), eventPayload)
		c.EventHandler.HandleEvent(event)
	case commands.DeleteAccount:
		account, err := c.AccountsService.GetAccountByAccountId(commandPayload.AccountId)
		if err != nil {
			return events.Event{}, err
		}

		event := events.NewEvent(time.Now().UTC().Round(time.Second), events.AccountDeleted{AccountId: account.AccountId})
		c.EventHandler.HandleEvent(event)
	case commands.CreateCommentThread:
		// Command not handled by Accounts
	case commands.CreateComment:
		// Command not handled by Accounts
	case commands.DeleteComment:
		// Command not handled by Accounts
	default:
		return events.Event{}, fmt.Errorf("unrecognized command type : %s", commandPayload.CommandType())
	}
	return events.Event{}, nil
}

func NewCommandHandler(db *sqlx.DB) CommandHandler {
	accountsService := &Accounts{DB: db}
	return CommandHandler{
		AccountsService: accountsService,
		EventHandler: &EventHandler{
			AccountsService: accountsService,
		},
	}
}

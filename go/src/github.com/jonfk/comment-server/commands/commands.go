package commands

import (
	"encoding/json"
	"fmt"
	_ "time"

	"github.com/satori/go.uuid"
)

const (
	CreateAccountTypeName       = "CreateAccount"
	DeleteAccountTypeName       = "DeleteAccount"
	LoginAccountTypeName        = "LoginAccount"
	CreateCommentThreadTypeName = "CreateCommentThread"
	CreateCommentTypeName       = "CreateComment"
	DeleteCommentTypeName       = "DeleteComment"
)

type Command struct {
	CommandType string         `json:"commandType"`
	Payload     CommandPayload `json:"payload"`
}

type CommandPayload interface {
	CommandType() string
}

type CreateAccount struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type DeleteAccount struct {
	AccountId uuid.UUID `json:"accountId"`
}

type LoginAccount struct {
	EmailOrUsername string `json:"emailOrUsername"`
	Password        string `json:"password"`
}

type CreateCommentThread struct {
	PageUrl string `json:"pageUrl,omitempty"`
	Title   string `json:"title,omitempty"`
}

type CreateComment struct {
	Data            string     `json:"data,omitempty"`
	ParentId        *uuid.UUID `json:"parentId,omitempty"`
	CommentThreadId uuid.UUID  `json:"commentThreadId,omitempty"`
	AccountId       uuid.UUID  `json:"accountId,omitempty"`
}

type DeleteComment struct {
	CommentId uuid.UUID `json:"commentId,omitempty"`
	AccountId uuid.UUID `json:"accountId,omitempty"`
}

func (c CreateAccount) CommandType() string       { return CreateAccountTypeName }
func (c DeleteAccount) CommandType() string       { return DeleteAccountTypeName }
func (c LoginAccount) CommandType() string        { return LoginAccountTypeName }
func (c CreateCommentThread) CommandType() string { return CreateCommentThreadTypeName }
func (c CreateComment) CommandType() string       { return CreateCommentTypeName }
func (c DeleteComment) CommandType() string       { return DeleteCommentTypeName }

type CommandJSON struct {
	CommandType string          `json:"commandType"`
	Payload     json.RawMessage `json:"payload"`
}

func UnmarshalJSON(input []byte) (CommandPayload, error) {
	var (
		commandRaw CommandJSON
	)
	err := json.Unmarshal(input, &commandRaw)
	if err != nil {
		return CreateAccount{}, err
	}

	switch commandRaw.CommandType {
	case CreateAccountTypeName:
		commandPayload := CreateAccount{}
		err = json.Unmarshal(commandRaw.Payload, &commandPayload)
		if err != nil {
			return commandPayload, err
		}
		return commandPayload, nil
	case DeleteAccountTypeName:
		commandPayload := DeleteAccount{}
		err = json.Unmarshal(commandRaw.Payload, &commandPayload)
		if err != nil {
			return commandPayload, err
		}
		return commandPayload, nil
	case LoginAccountTypeName:
		commandPayload := LoginAccount{}
		err = json.Unmarshal(commandRaw.Payload, &commandPayload)
		if err != nil {
			return commandPayload, err
		}
		return commandPayload, nil
	case CreateCommentThreadTypeName:
		commandPayload := CreateCommentThread{}
		err = json.Unmarshal(commandRaw.Payload, &commandPayload)
		if err != nil {
			return commandPayload, err
		}
		return commandPayload, nil
	case CreateCommentTypeName:
		commandPayload := CreateComment{}
		err = json.Unmarshal(commandRaw.Payload, &commandPayload)
		if err != nil {
			return commandPayload, err
		}
		return commandPayload, nil
	case DeleteCommentTypeName:
		commandPayload := DeleteComment{}
		err = json.Unmarshal(commandRaw.Payload, &commandPayload)
		if err != nil {
			return commandPayload, err
		}
		return commandPayload, nil
	default:
		return CreateAccount{}, fmt.Errorf("unknown command type %s", commandRaw.CommandType)
	}
}

func MarshalJSON(command Command) ([]byte, error) {
	return json.Marshal(command)
}

func CreateCommand(payload CommandPayload) Command {
	return Command{
		CommandType: payload.CommandType(),
		Payload:     payload,
	}
}

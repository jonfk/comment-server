package commands

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/satori/go.uuid"
)

func TestJsonMarshalCommands(t *testing.T) {
	parentId := uuid.NewV4()
	commandPayloads := []CommandPayload{
		CreateAccount{
			Username:         "username",
			Email:            "email",
			UnhashedPassword: "unhashed_password",
		},
		DeleteAccount{
			Username: "username",
			Email:    "email",
		},
		CreateCommentThread{
			PageUrl: "pageUrl",
			Title:   "title",
		},
		CreateComment{
			Data:            "this is data",
			ParentId:        &parentId,
			CommentThreadId: uuid.NewV4(),
			AccountId:       uuid.NewV4(),
		},
		CreateComment{
			Data:            "this is data",
			ParentId:        nil,
			CommentThreadId: uuid.NewV4(),
			AccountId:       uuid.NewV4(),
		},
		DeleteComment{
			CommentId: uuid.NewV4(),
			AccountId: uuid.NewV4(),
		},
	}

	encodedCommands := [][]byte{}
	decodedCommandPayloads := []CommandPayload{}

	for _, commandPayload := range commandPayloads {
		command := CreateCommand(commandPayload)
		encodedCommand, err := json.MarshalIndent(command, "", "  ")
		if err != nil {
			t.Fatalf("json.MarshalIndent failed : %v\n", err)
		}
		encodedCommands = append(encodedCommands, encodedCommand)
	}

	for _, encodedCommand := range encodedCommands {
		commandPayload, err := UnmarshalJSON(encodedCommand)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed : %v\n", err)
		}
		decodedCommandPayloads = append(decodedCommandPayloads, commandPayload)
	}

	for i, decodedCommandPayload := range decodedCommandPayloads {
		if !reflect.DeepEqual(commandPayloads[i], decodedCommandPayload) {
			t.Fatalf("commandPayload != decodedCommandPayload\n (commandPayload) %v != (decodedCommandPayload) %v", commandPayloads[i], decodedCommandPayload)
		}
	}

}

func TestUnmarshalJSON(t *testing.T) {
	jsonCommands := []string{
		`{
	"commandType": "CreateAccount",
	"payload": {
		"username": "username",
		"email": "email",
		"password": "password"
	}
}`,
		`{
	"commandType": "DeleteAccount",
	"payload": {
		"username": "username",
		"email": "email"
	}
}`,
		`{
	"commandType": "CreateCommentThread",
	"payload": {
		"pageUrl": "pageUrl",
		"title": "title"
	}
}`,
		`{
	"commandType": "CreateComment",
	"payload": {
		"data": "this is data",
		"parentId": "5a8433e3-4f28-4ad5-8f07-851e820b3205",
		"commentThreadId": "6c0904a0-901f-47a2-96f7-821fd6f800c1",
		"accountId": "bce2547a-08b9-47bb-84f2-f90b9673bc6a"
	}
}`,
		`{
	"commandType": "CreateComment",
	"payload": {
		"data": "this is data",
		"parentId": null,
		"commentThreadId": "6c0904a0-901f-47a2-96f7-821fd6f800c1",
		"accountId": "bce2547a-08b9-47bb-84f2-f90b9673bc6a"
	}
}`,
		`{
	"commandType": "CreateComment",
	"payload": {
		"data": "this is data",
		"commentThreadId": "a34fd176-15ca-43f0-9814-552f6b903723",
		"accountId": "e88ba639-4809-45d5-8f3c-1f3438d7a2a9"
	}
}`,
		`{
	"commandType": "DeleteComment",
	"payload": {
		"commentId": "1c1b0d99-5108-458a-a93d-131cc00c717c",
		"accountId": "81203854-0f6e-4b3e-88fd-4b6586cffa07"
	}
}`,
	}

	invalidJsonCommands := []string{
		`{
	"commandType": "CreateComment",
	"payload": {
		"data": "this is data",
		"parentId": "",
		"commentThreadId": "6c0904a0-901f-47a2-96f7-821fd6f800c1",
		"accountId": "bce2547a-08b9-47bb-84f2-f90b9673bc6a"
	}
}`,
	}

	for _, jsonCommand := range jsonCommands {
		_, err := UnmarshalJSON([]byte(jsonCommand))
		if err != nil {
			t.Fatalf("UnmarshalJSON failed on %s with %v\n", jsonCommand, err)
		}
	}

	for _, jsonCommand := range invalidJsonCommands {
		_, err := UnmarshalJSON([]byte(jsonCommand))
		if err == nil {
			t.Fatalf("UnmarshalJSON did fail as expected on %s\n", jsonCommand)
		}
	}

}


export GOPATH=$(shell pwd)

.PHONY: install clean get-deps test unit-test run

install:
	go install github.com/jonfk/comment-server
	GOOS=linux GOARCH=amd64 go install github.com/jonfk/comment-server

clean:
	rm -rf ./bin/
	rm -rf ./pkg/

get-deps:
	cd src/github.com/jonfk/comment-server && glide install

test:
	go test -v -cover github.com/jonfk/comment-server/accounts
	go test -v -cover github.com/jonfk/comment-server/comments
	go test -v -cover github.com/jonfk/comment-server/commands
	go test -v -cover github.com/jonfk/comment-server/events

unit-test:
	go test -v -short -cover github.com/jonfk/comment-server/accounts
	go test -v -short -cover github.com/jonfk/comment-server/comments
	go test -v -short -cover github.com/jonfk/comment-server/commands
	go test -v -short -cover github.com/jonfk/comment-server/events

run:
	# commands to run during development

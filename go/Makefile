
export GOPATH=$(shell pwd)

.PHONY: install clean get-deps test unit-test

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

unit-test:
	go test -v -short -cover github.com/jonfk/comment-server/accounts

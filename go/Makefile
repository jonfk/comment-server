
export GOPATH=$(shell pwd)

.PHONY: install clean get-deps

install:
	go install github.com/jonfk/comment-server

clean:
	rm -rf ./bin/
	rm -rf ./pkg/

get-deps:
	cd src/github.com/jonfk/comment-server && glide install

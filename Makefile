##=======================================================================##
## Makefile
## Created: Wed Aug 05 14:35:14 PDT 2015 @941 /Internet Time/
# :mode=makefile:tabSize=3:indentSize=3:
## Purpose: 
##======================================================================##

SHELL=/bin/bash
GPATH = $(shell pwd)

.PHONY: fmt install get scrape build clean 

install: fmt
	@GOPATH=${GPATH} go build -o skeleton-cli client.go

fmt:
	@GOPATH=${GPATH} gofmt -s -w skeleton_db
	@GOPATH=${GPATH} gofmt -s -w diff_store
	@GOPATH=${GPATH}  gofmt -s -w client.go

get:
	@GOPATH=${GPATH} go get ${OPTS} ${ARGS}

scrape:
	@find src -type d -name '.hg' -or -type d -name '.git' | xargs rm -rf

clean:
	@GOPATH=${GPATH} go clean

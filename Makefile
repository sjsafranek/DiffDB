##=======================================================================##
## Makefile
## Created: Wed Aug 05 14:35:14 PDT 2015 @941 /Internet Time/
# :mode=makefile:tabSize=3:indentSize=3:
## Purpose: 
##======================================================================##

SHELL=/bin/bash
GPATH = $(shell pwd)

.PHONY: fmt deps install build scrape clean

install: fmt deps
	@GOPATH=${GPATH} go build -o skeleton-cli client.go

build: fmt deps
	@GOPATH=${GPATH} go build -o skeleton-cli client.go

deps:
	mkdir -p "src"
	mkdir -p "pkg"
	@GOPATH=${GPATH} go get github.com/boltdb/bolt
	@GOPATH=${GPATH} go get github.com/sergi/go-diff/diffmatchpatch

fmt:
	@GOPATH=${GPATH} gofmt -s -w skeleton_db
	@GOPATH=${GPATH} gofmt -s -w diff_store
	@GOPATH=${GPATH} gofmt -s -w client.go

#get:
#	@GOPATH=${GPATH} go get ${OPTS} ${ARGS}

scrape:
	@find src -type d -name '.hg' -or -type d -name '.git' | xargs rm -rf

clean:
	@GOPATH=${GPATH} go clean

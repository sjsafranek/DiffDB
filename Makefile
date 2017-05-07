##=======================================================================##
## Makefile
## Created: Wed Aug 05 14:35:14 PDT 2015 @941 /Internet Time/
# :mode=makefile:tabSize=3:indentSize=3:
## Purpose:
##======================================================================##

SHELL=/bin/bash
PROJECT_NAME = DiffDB
GPATH = $(shell pwd)

.PHONY: fmt get-deps test install build scrape clean

install: fmt get-deps
	#@GOPATH=${GPATH} go build -o db-cli client.go

build: fmt get-deps
	#@GOPATH=${GPATH} go build -o db-cli client.go

get-deps:
	mkdir -p "src"
	mkdir -p "pkg"
	@GOPATH=${GPATH} go get github.com/boltdb/bolt
	@GOPATH=${GPATH} go get github.com/sergi/go-diff/diffmatchpatch
	@GOPATH=${GPATH} go get github.com/sjsafranek/SkeletonDB

fmt:
	@GOPATH=${GPATH} gofmt -s -w diff_db
	@GOPATH=${GPATH} gofmt -s -w diff_store
	#@GOPATH=${GPATH} gofmt -s -w client.go

test: fmt get-deps
	cd diff_store
	@GOPATH=${GPATH} go test -v -bench=. -test.benchmem

scrape:
	@find src -type d -name '.hg' -or -type d -name '.git' | xargs rm -rf

clean:
	@GOPATH=${GPATH} go clean

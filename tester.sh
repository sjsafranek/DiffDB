#!/bin/bash

export GOPATH="`pwd`"

cd src/ts_db
go test -bench=. -test.benchmem

#!/bin/bash

export GOPATH="`pwd`"

cd src/skeleton_db
go test -bench=. -test.benchmem

#!/bin/bash

export GOBIN="$(pwd)/bin"
export CGO_ENABLED=0

go build -tags netgo -a
go install



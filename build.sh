#!/bin/sh

PWD="$(pwd)"
ROOT="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

cd "$ROOT/cmd/exporter"

export CGO_ENABLED=0
go build -a -o "$ROOT/bin/prom-smartctl-exporter" -tags netgo

cd "$PWD"
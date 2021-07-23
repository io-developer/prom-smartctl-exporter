#!/bin/sh

docker run --rm -v $(pwd):/app golang:1.16 /app/build.sh

docker build --tag iodeveloper/prom-smartctl-exporter:latest .

#!/usr/bin/env bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -i -o http-test  http-test.go
docker build -t http-test:v0.1 .
#!/usr/bin/env bash

CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/weatherapp_amd64_linux ./main.go

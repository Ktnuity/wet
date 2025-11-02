#!/bin/bash

set -e

go mod tidy
go build -o wet cmd/cli/main.go

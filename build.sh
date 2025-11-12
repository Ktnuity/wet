#!/bin/bash

set -e

go mod tidy
go vet ./...
go generate ./...
go build -o wet cmd/cli/main.go
go build -o test cmd/test/main.go

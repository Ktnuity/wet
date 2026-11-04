#!/bin/bash

set -e

go generate ./...
go mod tidy
go vet ./...
go build -o wet cmd/cli/main.go
go build -o test cmd/test/main.go

rm -rf internal/stdlib/std

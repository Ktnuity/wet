#!/bin/bash

set -e

if [[ ! -d "./bin" ]]; then
    mkdir ./bin
fi

go generate ./...
go mod tidy
go vet ./...
go build -o bin/wet cmd/cli/main.go
go build -o bin/test cmd/test/main.go

rm -rf internal/stdlib/std

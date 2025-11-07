#!/bin/bash

set -e

go mod tidy
go generate cmd/cli/main.go
go build -o wet cmd/cli/main.go

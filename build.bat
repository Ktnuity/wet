go generate ./...
go mod tidy
go vet ./...
go build -o wet.exe cmd/cli/main.go
go build -o test.exe cmd/test/main.go

rmdir internal\stdlib\std


if not exist bin mkdir bin

go generate ./...
go mod tidy
go vet ./...
go build -o bin\wet.exe cmd/cli/main.go
go build -o bin\test.exe cmd/test/main.go

rmdir internal\stdlib\std


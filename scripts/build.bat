@echo off

if not exist bin mkdir bin

go generate ./...
if errorlevel 1 exit /b 1

go mod tidy
if errorlevel 1 exit /b 1

go vet ./...
if errorlevel 1 exit /b 1

go build -o bin\wet.exe cmd/cli/main.go
if errorlevel 1 exit /b 1

go build -o bin\test.exe cmd/test/main.go
if errorlevel 1 exit /b 1

if exist internal\stdlib\std rmdir /s /q internal\stdlib\std


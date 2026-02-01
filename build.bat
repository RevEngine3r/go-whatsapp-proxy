@echo off
SETLOCAL EnableDelayedExpansion

SET APP_NAME=whatsapp-proxy
SET BUILD_DIR=build

if not exist %BUILD_DIR% mkdir %BUILD_DIR%

echo Building for Windows x64...
SET GOOS=windows
SET GOARCH=amd64
go build -o %BUILD_DIR%/%APP_NAME%-windows-amd64.exe .

echo Building for Windows ARM64...
SET GOARCH=arm64
go build -o %BUILD_DIR%/%APP_NAME%-windows-arm64.exe .

echo Building for Linux x64...
SET GOOS=linux
SET GOARCH=amd64
go build -o %BUILD_DIR%/%APP_NAME%-linux-amd64 .

echo Building for Linux ARM64...
SET GOARCH=arm64
go build -o %BUILD_DIR%/%APP_NAME%-linux-arm64 .

echo Building for macOS x64...
SET GOOS=darwin
SET GOARCH=amd64
go build -o %BUILD_DIR%/%APP_NAME%-darwin-amd64 .

echo Building for macOS ARM64...
SET GOARCH=arm64
go build -o %BUILD_DIR%/%APP_NAME%-darwin-arm64 .

echo Done! Binaries are in the %BUILD_DIR% folder.
pause

#!/bin/sh

PACKAGE=./cmd/slimiter-server/
RELEASE_PATH=release/slimiter-server

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 vgo build -o $RELEASE_PATH-darwin $PACKAGE
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 vgo build -o $RELEASE_PATH-linux $PACKAGE
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 vgo build -o $RELEASE_PATH-freebsd $PACKAGE
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 vgo build -o $RELEASE_PATH-windows.exe $PACKAGE

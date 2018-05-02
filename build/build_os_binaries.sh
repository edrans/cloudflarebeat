#!/bin/bash

cd /go/src/$WORKDIR
go get -d -v ./...
go install -v ./...


mkdir bin
echo "Building v$1 for Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/cloudflarebeat-$VERSION-linux-amd64
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/cloudflarebeat-$VERSION-darwin-amd64
#GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/cloudflarebeat-$1-windows-amd64

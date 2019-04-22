#!/bin/sh -xe

rm -fr dist && mkdir dist

mkdir -p build/amd64
GOOS=linux GOARCH=amd64 go build -o build/amd64/auto-builder cmd/auto-builderd/main.go

mkdir -p pkg-build/amd64
go-bin-deb generate -a amd64 -w pkg-build/amd64/ --version 0.1.0 -o dist/auto-builder-amd64.deb

rm -fr build
rm -fr pkg-build

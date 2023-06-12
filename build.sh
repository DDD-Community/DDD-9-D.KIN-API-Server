#!/bin/bash

rm -r build/
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o build/main main.go
cd build
zip function.zip main
rm main
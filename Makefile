MODULE_NAME := $(shell head -n 1 go.mod | awk '{print $2}')
GO_BIN_PATH := $(shell go env GOPATH)/bin
PROJECT_DIR := $(shell pwd)

LAMLAM_CMD := $(GO_BIN_PATH)/lamlam

install:
	go mod download all
	go install github.com/stockfolioofficial/lamlam/cmd/lamlam@v1.0.0

install-oas-mac:
	brew install openapi-generator

oas:
	cd $(PROJECT_DIR)/docs/oas && openapi-generator generate -i docs.yaml -g openapi-yaml -o gen

gen: oas

lambda-build:
	-rm -r build/
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o build/main main.go
	cd build && zip function.zip main && rm main

lambda-deploy:
	go run ./cmd/deploy

.PHONY: create-env-file install install-oas-mac oas

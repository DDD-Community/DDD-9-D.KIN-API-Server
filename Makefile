MODULE_NAME := $(shell head -n 1 go.mod | awk '{print $2}')
GO_BIN_PATH := $(shell go env GOPATH)/bin
PROJECT_DIR := $(shell pwd)

LAMLAM_CMD := $(GO_BIN_PATH)/lamlam


define ENV_FORMAT
AWS_REGION=ap-northeast-2
AWS_ACCESS_KEY=INPUT_YOUR_ACCESS_KEY
AWS_SECRET_ACCESS_KEY=INPUT_YOUR_SECRET_ACCESS_KEY
endef
export ENV_FORMAT

local-env-file:
	echo "$$ENV_FORMAT" > .env.local

install:
	go mod download all
	go install github.com/stockfolioofficial/lamlam/cmd/lamlam@v1.0.0

install-oas-mac:
	brew install openapi-generator

oas:
	cd $(PROJECT_DIR)/docs/oas && openapi-generator generate -i docs.yaml -g openapi-yaml -o gen

gen: oas

# TODO: 우선 묶어서 전체 배포
lambda-build:
	-rm -r build/
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o build/app main.go
	cd build && zip app.zip app && rm app

	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o build/authorizer cmd/lambda/authorizer/main.go
	cd build && zip authorizer.zip authorizer && rm authorizer

lambda-deploy:
	go run ./cmd/deploy

.PHONY: local-env-file install install-oas-mac oas gen lambda-build lambda-deploy

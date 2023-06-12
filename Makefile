ifndef RUNTIME
	RUNTIME=local
endif

ENV_FILE := .env.$(RUNTIME)

ifeq (,$(wildcard ./$(ENV_FILE)))
    ENV_FILE_EXISTS=0
else
	ENV_FILE_EXISTS=1
	include $(ENV_FILE)
endif

MODULE_NAME := $(shell head -n 1 go.mod | awk '{print $2}')
GO_BIN_PATH := $(shell go env GOPATH)/bin
PROJECT_DIR := $(shell pwd)

LAMLAM_CMD := $(GO_BIN_PATH)/lamlam

define ENV_FORMAT
AWS_REGION=ap-northeast-2
AWS_ACCESS_KEY=INPUT_YOUR_ACCESS_KEY
AWS_SECRET_ACCESS_KEY=INPUT_YOUR_SECRET_ACCESS_KEY
SERVER_PORT=3000
endef
export ENV_FORMAT

INPUT_NEW_ENV_FILE ?= $(shell bash -c 'read -p "Write env stage (default. local): " env_file_name; echo $$env_file_name')

create-env-file:
	@NEW_ENV_FILE=$(INPUT_NEW_ENV_FILE); \
	if [ "$(NEW_ENV_FILE)" != "" ]; \
		then \
		  echo "$$ENV_FORMAT" > .env.$(NEW_ENV_FILE); \
		else \
		  echo "$$ENV_FORMAT" > .env.local; \
  	fi

install:
	go mod download all
	go install github.com/stockfolioofficial/lamlam/cmd/lamlam@v1.0.0

install-oas-mac:
	brew install openapi-generator

oas:
	cd $(PROJECT_DIR)/docs && openapi-generator generate -i docs.yaml -g openapi-yaml -o gen

gen: oas

.PHONY: create-env-file install install-oas-mac oas

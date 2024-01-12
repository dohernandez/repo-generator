GO ?= go

PWD = $(shell pwd)

# Detecting GOPATH and removing trailing "/" if any
GOPATH = $(realpath $(shell $(GO) env GOPATH))

TEST_REPO_GENERATOR_DEVGO_PATH ?= $(PWD)/internal/makefiles
TEST_REPO_GENERATOR_DEVGO_SCRIPTS ?= $(TEST_REPO_GENERATOR_DEVGO_PATH)

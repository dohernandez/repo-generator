GO ?= go

## Check/install repo-generator tool
repo-generator-cli:
	@REPO_GENERATOR_VERSION=$(REPO_GENERATOR_VERSION) bash $(REPO_GENERATOR_DEVGO_SCRIPTS)/repo-generator-cli.sh

.PHONY: repo-generator-cli

#-## Create/replace GitHub Actions from template

#- Placeholders require include the file in the Makefile
#- require - dev/github-actions

GO ?= go

AFTER_GITHUB_TARGETS += github-dependabot-testdata

## Inject/Replace GitHub dependabot test repo-generator output
github-dependabot-testdata:
	@echo "Updating dependabot.yml"
	@bash $(TEST_REPO_GENERATOR_DEVGO_SCRIPTS)/github.sh

.PHONY: github-dependabot-testdata

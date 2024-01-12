#!/bin/bash

# The path to the dependabot file
dependabot_path=".github/dependabot.yml"

# Add the dependabot gomod testdata
yq e '.updates += [
  {
    "package-ecosystem": "gomod",
    "directory": "/testdata",
    "schedule": {"interval": "daily"},
    "open-pull-requests-limit": 1,
    "groups": {
      "test-dependencies": {
        "patterns": ["*"],
        "update-types": ["minor", "patch"]
      }
    }
  }
]' -i $dependabot_path
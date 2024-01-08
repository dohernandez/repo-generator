#!/usr/bin/env bash

[ -z "$REPO_GENERATOR_VERSION" ] && REPO_GENERATOR_VERSION="0.2.1"

install_repo_generator () {
  case "$2" in
      Darwin*)
        {
          PLATFORM=$2
          HARDWARE=$(uname -m)
        };;
      Linux*)
        {
          PLATFORM=$2
          HARDWARE=$(uname -m)
           if [ "$HARDWARE" == "aarch64" ]; then \
                HARDWARE="arm64"
              fi
        };;
  esac

  mkdir -p /tmp/repo_generator-"$1"

  curl -sL https://github.com/dohernandez/repo-generator/releases/download/v"$1"/"$PLATFORM"_"$HARDWARE".tar.gz | tar xvz -C /tmp/repo_generator-"$1" \
    && mv /tmp/repo_generator-"$1"/repo-generator "$GOPATH"/bin/repo-generator
}

osType="$(uname -s)"

# detecting GOPATH and removing trailing "/" if any
GOPATH="$(go env GOPATH)"
GOPATH=${GOPATH%/}

# adding GOBIN to PATH
[[ ":$PATH:" != *"$GOPATH/bin"* ]] && PATH=$PATH:"$GOPATH"/bin

# checking if repo-generator is available and it is the version specify
if ! command -v repo-generator > /dev/null; then \
    echo ">> Installing repo-generator v$REPO_GENERATOR_VERSION...";
    install_repo_generator "$REPO_GENERATOR_VERSION" "$osType"
else
  VERSION_INSTALLED="$(repo-generator --version --quiet | cut -d' ' -f2)"
  if [ "${VERSION_INSTALLED}" != "v${REPO_GENERATOR_VERSION}" ]; then \
    echo ">> Updating repo-generator form "${VERSION_INSTALLED}" to v$REPO_GENERATOR_VERSION..."; \
    install_repo_generator "$REPO_GENERATOR_VERSION" "$osType"
  fi
fi

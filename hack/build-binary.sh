#!/usr/bin/env bash

set -o errexit
set -o nounset

__usage() {
  cat <<EOF
USAGE:

${0} [BINARY_NAME]

Required environment variables:
    GO_BUILD_LDFLAGS    go linker flags.
EOF
  exit 1
}

[ -z "${1}" ] && __usage
[ -z "${GO_BUILD_LDFLAGS}" ] && __usage

BINARY_NAME="${1}"

export CGO_ENABLED=0

go build \
  -ldflags "${GO_BUILD_LDFLAGS}" \
  -o "build/bin/${BINARY_NAME}" \
  "./cmd/${BINARY_NAME}"

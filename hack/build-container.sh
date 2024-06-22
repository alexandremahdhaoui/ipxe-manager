#!/usr/bin/env bash

set -xe

__usage() {
  cat <<EOF
USAGE:

${0} [BINARY_NAME]

Required environment variables:
    CONTAINER_ENGINER   container engine such as podman or docker.
    GO_BUILD_LDFLAGS    go linker flags.
    VERSION             tag semver.
EOF
  exit 1
}

[ -z "${1}" ] && __usage
[ -z "${CONTAINER_ENGINE}" ] && __usage
[ -z "${GO_BUILD_LDFLAGS}" ] && __usage

BINARY_NAME="${1}"

"${CONTAINER_ENGINE}" \
  build \
  . \
  --build-arg "GO_BUILD_LDFLAGS=${GO_BUILD_LDFLAGS}" \
  -t "${BINARY_NAME}:${VERSION}" \
  -f "./containers/${BINARY_NAME}/Containerfile"

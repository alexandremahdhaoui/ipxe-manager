#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

function __usage() {
  cat <<EOF
${0} [CMD]

Available commands:
  setup       This command will set up the e2e test environment.
  run         This command will run the e2e tests.
  teardown    This command will tear down the e2e test environment.

  full-test   This command will set up the environment, run the 
              end-to-end tests, and tear the environment down.
EOF
}

BRIDGE_IFACE=e2e0br0
DUMMY_IFACE=e2edummy0
TAP_IFACE=e2e0tap0

WDIR="$(git rev-parse --show-toplevel)"
TEMPDIR="${WDIR}/.tmp/e2e"

DNSMASQ_PID_FILE="${TEMPDIR}/dnsmasq.pid"
DNSMASQ_LOG="${TEMPDIR}/dnsmasq.log"
DNSMASQ_PROCESS_LOG="${TEMPDIR}/dnsmasq.process.log"
DNSMASQ_TFTP_DIR="${TEMPDIR}/tftpboot"
DNSMASQ_CONF_FILE="${TEMPDIR}/dnsmasq.conf"

export BRIDGE_IFACE DNSMASQ_LOG DNSMASQ_TFTP_DIR DNSMASQ_PID_FILE

function __setup() {
  echo "⏳ Setting up e2e environment..."

  # -- Create e2e temp directory.
  mkdir -p "${DNSMASQ_TFTP_DIR}"

  # -- Generate dnsmasq config
  envsubst <"${WDIR}/e2e/templates/dnsmasq.conf.tmpl" | tee "${DNSMASQ_CONF_FILE}" 1>/dev/null

  # -- Create bridge interface.
  sudo ip l add dev "${BRIDGE_IFACE}" type bridge
  sudo ip a add 172.16.0.1/24 dev "${BRIDGE_IFACE}"
  sudo ip l set dev "${BRIDGE_IFACE}" up

  # -- Create dummy interface.
  sudo ip link add name "${DUMMY_IFACE}" up master "${BRIDGE_IFACE}" type dummy

  # -- Create tap device for client.
  sudo ip tuntap add "${TAP_IFACE}" mode tap
  sudo ip l set dev "${TAP_IFACE}" master "${BRIDGE_IFACE}"
  sudo ip l set dev "${TAP_IFACE}" up

  # -- Run dnsmasq.
  echo "⏳ Starting dhcp server..."
  touch "${DNSMASQ_LOG}"
  sudo dnsmasq -d --conf-file="${DNSMASQ_CONF_FILE}" &>"${DNSMASQ_PROCESS_LOG}" &
  echo -n $! | tee "${DNSMASQ_PID_FILE}" 1>/dev/null

  echo "✅ Successfully set up e2e environment!"
}

function __run() {
  echo "TODO: run command"
  sudo dhclient -v -s 172.16.0.255 "${BRIDGE_IFACE}"
}

function __teardown() {
  set +o errexit

  echo "⏳ Tearing down e2e environment..."

  echo "⏳ Terminating dhcp server..."
  sudo kill -9 "$(cat "${DNSMASQ_PID_FILE}")"
  rm "${DNSMASQ_PID_FILE}"

  echo "⏳ Deleting network interfaces \"${BRIDGE_IFACE}\"..."
  sudo ip l del dev "${BRIDGE_IFACE}"
  sudo ip l del dev "${DUMMY_IFACE}"
  sudo ip l del dev "${TAP_IFACE}"

  echo "✅ Successfully deleted e2e environment!"
  set -o errexit
}

trap usage EXIT
CMD="${1}"
trap : EXIT

function main() {
  case "${CMD}" in
  setup)
    __setup
    exit 0
    ;;

  run)
    __run
    exit 0
    ;;

  teardown)
    __teardown
    exit 0
    ;;

  full-test)
    trap __teardown EXIT
    __setup
    __run

    trap : EXIT
    __teardown
    ;;

  *)
    __usage
    exit 1
    ;;
  esac
}

main

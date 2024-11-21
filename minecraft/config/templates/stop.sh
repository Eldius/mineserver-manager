#!/bin/bash

handle_error() {
  msg=$1

  echo "error: $msg"
  exit 1
}

INSTALL_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

cd "${INSTALL_PATH}" || exit 1

PID="$( cat ${INSTALL_PATH}/server.pid )"

echo "stopping server process: $PID"

kill $PID && echo '' > ${INSTALL_PATH}/server.pid || handle_error "killing process"

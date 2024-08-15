#!/bin/bash

INSTALL_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

cd "${INSTALL_PATH}" || exit 1

"${INSTALL_PATH}/java/jdk/bin/java" -Xms{{ .Xms }} -Xmx{{ .Xmx }} {{ if .LogConfigFile }}-Dlog4j.configurationFile=${INSTALL_PATH}/log4j2.xml {{ end }} -jar server.jar {{ if .Headless }}--nogui{{ end }} &

PID=$!
echo $PID > server.pid
echo "starting server with PID: $PID"

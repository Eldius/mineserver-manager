#!/bin/bash

INSTALL_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

cd "${INSTALL_PATH}" || exit 1

"${INSTALL_PATH}/jdk-17.0.7+7/bin/java" -Xms{{ .Xms }} -Xmx{{ .Xmx }} {{ if .LogConfigFile }}-Dlog4j.configurationFile={{ .LogConfigFile }}{{ end }} -jar server.jar {{ if .Headless }}--nogui{{ end }}

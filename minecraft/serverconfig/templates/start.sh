#!/bin/bash

./jdk-17.0.7+7/bin/java -Xms{{ .Xms }} -Xmx{{ .Xmx }} {{ if .LogConfigFile }}-Dlog4j.configurationFile={{ .LogConfigFile }}{{ end }} -jar server.jar {{ if .Headless }}--nogui{{ end }}

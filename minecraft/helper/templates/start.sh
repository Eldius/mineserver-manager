#!/bin/bash

## Based on Aikar's script generator
## https://docs.papermc.io/paper/aikars-flags

INSTALL_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

cd "${INSTALL_PATH}" || exit 1

{{ .JDKPath }}/java {{ if .LogConfigFile }}-Dlog4j.configurationFile=${INSTALL_PATH}/log4j2.xml {{ end }} \
  -Xms{{ .MemLimit }} \
  -Xmx{{ .MemLimit }} \
  -XX:+UseG1GC \
  -XX:+ParallelRefProcEnabled \
  -XX:MaxGCPauseMillis=200 \
  -XX:+UnlockExperimentalVMOptions \
  -XX:+DisableExplicitGC \
  -XX:+AlwaysPreTouch \
  -XX:G1NewSizePercent=30 \
  -XX:G1MaxNewSizePercent=40 \
  -XX:G1HeapRegionSize=8M \
  -XX:G1ReservePercent=20 \
  -XX:G1HeapWastePercent=5 \
  -XX:G1MixedGCCountTarget=4 \
  -XX:InitiatingHeapOccupancyPercent=15 \
  -XX:G1MixedGCLiveThresholdPercent=90 \
  -XX:G1RSetUpdatingPauseTimePercent=5 \
  -XX:SurvivorRatio=32 \
  -XX:+PerfDisableSharedMem \
  -XX:MaxTenuringThreshold=1 \
  -Dusing.aikars.flags=https://mcflags.emc.gs \
  -Daikars.new.flags=true \
  -jar {{ .ServerFile }} {{ if .Headless }} --nogui {{ end }} &

PID=$!
echo $PID > {{ .PIDFile }}
echo "starting server with PID: $PID"

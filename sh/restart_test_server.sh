#!/bin/bash
scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir=$(echo "${scriptdir%/*}")

ps aux | grep ".go" | grep "log-level trace" |
  awk '{print $2}' | xargs -I{} kill -15 {}

r ${basedir}/testdata/conf.yaml --log-level trace

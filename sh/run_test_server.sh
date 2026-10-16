#!/bin/bash
scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir=$(echo "${scriptdir%/*}")

ll="trace"
[[ -n "${1}" ]] && ll="${1}"

ps aux | grep ".go" | grep "log-level trace" |
  awk '{print $2}' | xargs -I{} kill -15 {}

r ${basedir}/testdata/conf.yaml --log-level "${ll}"

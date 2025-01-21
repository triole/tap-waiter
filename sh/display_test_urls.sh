#!/bin/bash
scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir=$(echo "${scriptdir%/*}")

cd "${basedir}/src/indexer" && go test >/dev/stdout 2>&1 |
  grep -Po "(?<=urls:\[).*?(?=\])" | sort

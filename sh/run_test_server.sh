#!/bin/bash
scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir=$(echo "${scriptdir%/*}")

ll="trace"
[[ -n "${1}" ]] && ll="${1}"

r ${basedir}/testdata/conf.yaml --log-level "${ll}"

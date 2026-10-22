#!/bin/bash

BASE_URL="http://localhost:17777"
ENDPOINT="/all.json"

fetch() {
  echo -e "\n$(date "+%Y-%m-%d %H:%M:%S") \033[0;95m${1}\033[0m"
  curl -s "${BASE_URL}${ENDPOINT}?${2}" | jq "${3}"
}

fetch "size ascending order" "sortby=size&order=asc" ".[0:5][].size"

fetch "size descending order" "sortby=size&order=desc" ".[0:5][].size"

fetch "filter tags contain tag1" "filter=front_matter.tags==tag1" \
  ".[].content.front_matter.tags"

fetch "filter tags equal tag1" "filter=front_matter.tags===tag1" \
  ".[].content.front_matter.tags"

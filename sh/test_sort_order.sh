#!/bin/bash

BASE_URL="http://localhost:17777"
ENDPOINT="/all.json"

echo "size ascending order"
curl -s "${BASE_URL}${ENDPOINT}?sortby=size&order=asc" | jq '.[0:3][] | {path, size}'

echo "size descending order"
curl -s "${BASE_URL}${ENDPOINT}?sortby=size&order=desc" | jq '.[0:3][] | {path, size}'

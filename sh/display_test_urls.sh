#!/bin/bash

task test >/dev/stdout 2>&1 |
  grep -Po "(?<=url\=).*" |
  grep -Po "\?.*" | sort

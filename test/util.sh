#!/usr/bin/env bash

# this will trap any errors or commands with non-zero exit status
# by calling function catch_errors()
trap catch_errors ERR;

check_success() {
  if [ -z "$1" ] || [ "$1" = "null" ]; then
    echo "Error: $2 is null"
    exit 1
  elif $(grep -q "error" <<< "$1"); then
    echo "Error: $2"
    exit 1
  else
    echo "Success!!"
  fi
}

assert_equals() {
  if [ "$1" != "$2" ]; then
    echo "Error: $3"
    exit 1
  else
    echo "Success!!"
  fi
}
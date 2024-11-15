#!/usr/bin/env bash

check_success() {
  if [ -z "$1" ] || [ "$1" = "null" ]; then
    echo "Error: $2 is null"
    exit 1
  elif $(grep -q "error" <<< "$1"); then
    echo "Error: $2 has an error"
    echo "Value: $1"
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
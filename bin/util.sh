#!/bin/bash

# this will trap any errors or commands with non-zero exit status
# by calling function catch_errors()
trap catch_errors ERR;

function catch_errors() {
   exit_code=$?
   echo "### FAILURE ###";
   exit $exit_code;
}

# Check whether given command exists and call it with the --version flag.
check_command_exists() {
  if sudo-linux /bin/bash -c "command -v "$1" > /dev/null 2>&1"; then
    echo "$1 version: $(sudo-linux "$1" --version)"
  else
    echo "RADAR Platform cannot start without $1. Please, install $1 and then try again"
    exit 1
  fi
}

# Check if the parent directory of given variable is set. Usage:
# check_parent_exists MY_PATH_VAR $MY_PATH_VAR
check_parent_exists() {
  if [ -z "$2" ]; then
    echo "Directory variable $1 is not set in .env"
  fi
  PARENT=$(dirname $2)
  if [ ! -d "${PARENT}" ]; then
    echo "RADAR-base stores volumes at ${PARENT}. If this folder does not exist, please create the entire path and then try again"
    exit 1
  fi
  if [ -d "$2" ]; then
    sudo-linux chmod 700 "$2"
  else
    sudo-linux mkdir -p -m 0700 "$2"
  fi
}

# sudo if on Linux, not on OS X
# useful for docker, which doesn't need sudo on OS X
sudo-linux() {
  if [ $(uname) == "Darwin" ] || id -Gn | grep -qve '\<sudo\>'; then
    "$@"
  else
    sudo "$@"
  fi
}

# OS X/linux portable sed -i
sed_i() {
  if [ $(uname) == "Darwin" ]; then
    sed -i '' "$@"
  else
    sed -i -- "$@"
  fi
}

# Inline variable into a file, keeping indentation.
# Usage:
# inline_variable VARIABLE_SET VALUE FILE
# where VARIABLE_SET is a regex of the pattern currently used in given file to set a variable to a value.
# Example:
# inline_variable 'a=' 123 test.txt
# will replace a line '  a=232 ' with '  a=123'
inline_variable() {
  sed_i 's|^\([[:space:]]*'"$1"'\).*$|\1'"$2"'|' "$3"
}

ensure_variable() {
  if grep -q "$1" "$3"; then
    inline_variable "$@"
  else
    echo "$1$2" >> "$3"
  fi
}

# Copies the template (defined by the given config file with suffix
# ".template") to intended configuration file, if the file does not
# yet exist.
copy_template_if_absent() {
  template=${2:-${1}.template}
  if [ ! -f "$1" ]; then
    if [ -e "$1" ]; then
      echo "Configuration file ${1} is not a regular file."
      exit 1
    else
      echo "Copying configuration file ${1} from template ${template}"
      sudo-linux cp -p "${template}" "$1"
    fi
  elif [ "$1" -ot "${template}" ]; then
    echo "Configuration file ${1} is older than its template"
    echo "${template}. Please edit ${1}"
    echo "to ensure it matches the template, remove it or run touch on it."
    exit 1
  fi
}

# Copies the template (defined by the given config file with suffix
# ".template") to intended configuration file.
copy_template() {
  template=${2:-${1}.template}
  if [ ! -f "$1" ] && [ -e "$1" ]; then
    echo "Configuration file ${1} is not a regular file."
    exit 1
  else
    echo "Copying configuration file ${1} from template ${template}"
    sudo-linux cp -p "${template}" "$1"
  fi
}

check_config_present() {
  template=${2:-${1}.template}
  if [ ! -f "$1" ]; then
    if [ -e "$1" ]; then
      echo "Configuration file ${1} is not a regular file."
    else
      echo "Configuration file ${1} is not present."
      echo "Please copy it from ${template} and modify it as needed."
    fi
    exit 1
  elif [ "$1" -ot "${template}" ]; then
    echo "Configuration file ${1} is older than its template ${template}."
    echo "Please edit ${1} to ensure it matches the template or run touch on it."
    exit 1
  fi
}

query_password() {
  echo $2
  stty -echo
  printf "Password: "
  read PASSWORD
  stty echo
  printf "\n"
  eval "$1=$PASSWORD"
}

ensure_env_default() {
  if eval '[ -z $'$1' ]'; then
    ensure_variable "$1=" "$2" .env
    eval "$1=$2"
  fi
}

ensure_env_password() {
  if eval '[ -z $'$1' ]'; then
    query_password $1 "$2"
    ensure_variable "$1=" "$PASSWORD" .env
  fi
}

#!/usr/bin/env bash

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
  version_flag=${2:---version}
  if /usr/bin/env bash -c "command -v "$1" > /dev/null 2>&1"; then
    echo "$1 version: $("$1" $version_flag)"
  else
    echo "RADAR Platform cannot start without $1. Please, install $1 and then try again"
    exit 1
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
      cp -p "${template}" "$1"
    fi
  elif [ "$1" -ot "${template}" ]; then
    echo "Configuration file ${1} is older than its template ${template}."
    echo "Please edit ${1} to ensure it matches the template, remove it or run "
    echo
    echo "  touch $1"
    echo
    echo "and try again."
    exit 1
  fi
}

create_production_yaml() {
  copy_template_if_absent etc/production.yaml etc/base.yaml
  sed -i "/_chart_version/d" etc/production.yaml
}

# Copies the template (defined by the given config file with suffix
# ".template") to intended configuration file.
copy_template() {
  template=${2:-${1}.template}
  if [ ! -f "$1" ] && [ -e "$1" ]; then
    echo "Configuration file ${1} is not a regular file."
    exit 1
  fi
  echo "Copying configuration file ${1} from template ${template}"
  cp -p "${template}" "$1"
}

generate_secret() {
  size=${1:-32}
  openssl rand -base64 $size | tr -d '+/'
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

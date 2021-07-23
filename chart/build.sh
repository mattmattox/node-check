#!/bin/bash
CWD=`pwd`
timestamp() {
  date "+%Y-%m-%d %H:%M:%S"
}
techo() {
  echo "$(timestamp): $*"
}
decho() {
  if [[ ! -z $DEBUG ]]
  then
    techo "$*"
  fi
}
setup-ssh() {
  echo "Setting up SSH key..."
  if [[ -z $SSH_KEY ]]
  then
    echo "SSH Key is missing"
    exit 1
  fi
  mkdir /root/.ssh && echo "$SSH_KEY" > /root/.ssh/id_rsa && chmod 0600 /root/.ssh/id_rsa
}

if [[ -z $Action ]]
then
  echo "Action must be set"
  exit 0
fi
setup-ssh

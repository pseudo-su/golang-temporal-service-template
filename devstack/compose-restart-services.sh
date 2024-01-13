#!/usr/bin/env bash

set -eo pipefail

log() { echo "$@" 1>&2; }

DEFAULT_COMPOSE_TOOL="podman compose"
DEVSTACK_COMPOSE_TOOL="${DEVSTACK_COMPOSE_TOOL:-$DEFAULT_COMPOSE_TOOL}"

function compose_restart_services () {
  declare filter="$1"
  if [[ "$filter" == "" ]]; then
    services=$($DEVSTACK_COMPOSE_TOOL -f ./devstack/compose.yaml ps --services);
  else
    services=$($DEVSTACK_COMPOSE_TOOL -f ./devstack/compose.yaml ps --services | grep $filter);
  fi

  $DEVSTACK_COMPOSE_TOOL -f ./devstack/compose.yaml restart $services
}

compose_restart_services $@

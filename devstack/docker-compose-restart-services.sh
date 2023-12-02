#!/usr/bin/env bash

set -eo pipefail

log() { echo "$@" 1>&2; }

function docker_compose_restart_services () {
  declare filter="$1"
  if [[ "$filter" == "" ]]; then
    services=$(docker-compose -f ./devstack/docker-compose.yaml ps --services);
  else
    services=$(docker-compose -f ./devstack/docker-compose.yaml ps --services | grep $filter);
  fi

  docker-compose -f ./devstack/docker-compose.yaml restart $services
}

docker_compose_restart_services $@

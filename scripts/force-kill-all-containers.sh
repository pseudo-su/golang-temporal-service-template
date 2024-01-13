#!/bin/bash

set -eo pipefail

log() { echo "$@" 1>&2; }

running_containers=`podman ps -q`
if [[ "$running_containers" ]]; then
  log "Stopping containers: $running_containers";
  podman kill $running_containers
fi

containers=`podman ps -a -q`
if [[ "$containers" ]]; then
  log "Removing containers $containers";
  podman rm $containers
fi

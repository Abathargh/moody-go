#!/bin/bash

#################################################################
#   This script automatically runs docker-compose using         #
#   the right dockerfiles for your CPU architecture.            #
#   If you don't know what to do just run this script.          #
#################################################################

ARCH=$(uname -m)
SUPPORTED=true

case $ARCH in
  x86_64)
    export DOCKERFILE_ARCH=Dockerfile
    ;;
  armv7l)
    export DOCKERFILE_ARCH=Dockerfile.arm32
    ;;
  aarch64)
    export DOCKERFILE_ARCH=Dockerfile.arm64
    ;;
  *)
    SUPPORTED=false
esac

if [[ "$SUPPORTED" == "true" ]]; then
  docker-compose up
else
  echo "CPU architecture not supported"
fi
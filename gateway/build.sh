#!/bin/bash
CMD=./cmd/
BUILD_TARGETS=$(find $CMD -maxdepth 2 -mindepth 2)

for target in $BUILD_TARGETS
do
  go build "$target"
done
#!/bin/bash

if [[ $# < 1 ]]; then
  echo "Choose one of ci|dev as the build target, exit ..."
  exit 1
fi

target="$1"

if [[ "${target}" != "ci" && "${target}" != "dev" ]]; then
  echo "Choose one of ci|dev as the build target, exit ..."
  exit 1
fi

docker login -u tricorderobservability
make -C devops/dev_image "build_and_push_${target}_image"

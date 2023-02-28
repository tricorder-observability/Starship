#!/bin/bash -ex

REGISTRY=docker.io/tricorderobservability

image="${REGISTRY}/base_build_image:v0.1"
ToT=$(git rev-parse --show-toplevel)

echo "Logging into docker ..."
docker login --username tricorderobservability

docker build ${ToT}/devops/base_build_image/ -t ${image}
docker push ${image}

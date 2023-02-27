#!/bin/bash -ex

REGISTRY=docker.io/tricorderobservability

image="${REGISTRY}/base_build_image:v0.1"
docker build . -t ${image}
docker push ${image}

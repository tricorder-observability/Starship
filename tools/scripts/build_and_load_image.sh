#!/bin/bash

if [[ $# -lt 1 ]]; then
  echo "go_image() container_image() automatically include a copy_[name]_tar"
  echo "rule to copy the output .tar file to bazel-bin path."
  echo
  echo "This is necessary because the implicit .tar target's output path"
  echo "is not stable."
  echo
  echo "See bazel/container_image.bzl for more details"
  echo
  echo "Need at least 1 argument, which is the image build rule " \
       "(not the .tar rule), exit ..."
  exit 1
fi

label=${1##*:}
path=${1%:*}
bazel_rule="${path}:copy_${label}_tar"

echo "Building ${bazel_rule} ..."
tar_file=$(bazel build "${bazel_rule}" 2>&1 | tee |  grep 'bazel-bin\/.*\/.*\.tar')

echo "Loading ${tar_file} ..."
docker load -i ${tar_file}

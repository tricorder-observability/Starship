#!/bin/bash -e

if [[ ! -f Dockerfile ]]; then
  echo "Could not find Dockerfile. Must be running in the same directory "
  echo "of the clang deb building container, exit ..."
  exit 1
fi

# This script builds LLVM from source inside a docker container, and package
# the build artifacts into a .deb package.

docker_image_tag="clang_deb:tricorder"
deb_name="clang-14_0-tricorder-0.deb"

# Build the docker image. The docker build command fetches clang & llvm sources
# and build them using clang's own instructions.
docker build . --tag ${docker_image_tag}

# Run docker image. The container invokes fpm to package the build artifacts
# into deb package.
# NOTE: This process is very slow.
docker run --interactive --tty --rm --env DEB_NAME=${deb_name} \
  --volume "$(pwd):/image" ${docker_image_tag}

s3_path="s3://tricorder-dev/${deb_name}"
echo "Uploading DEB package to ${s3_path} ..."
aws s3 cp ${deb_name} ${s3_path}

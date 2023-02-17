#!/bin/bash -ex

# Look up the build rule of a file in llvm Bazel repo:
# https://github.com/llvm/llvm-project/tree/main/utils/bazel
# This is needed when we were trying to include llvm project
# directly inside starship repo as a first-class external Bazel
# repo.
#
# The objective was not accomplished, because we haven't
# managed to replace the bazel CC toolchain with the clang
# compiler produced from this LLVM project, thus causes mysterious
# crash when BCC links this native llvm repo, and some other code
# links with the host's version.
#
# TODO(yaxiong): Use bazel managed llvm as cc toolchain.

echo $1
fullname=$(bazel query @llvm-project//llvm:$1)
echo ${fullname}
bazel query "attr('srcs', $fullname, ${fullname//:*/}:*)"
bazel query "attr('hdrs', $fullname, ${fullname//:*/}:*)"

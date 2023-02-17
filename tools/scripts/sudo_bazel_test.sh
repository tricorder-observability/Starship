#!/bin/bash -ex

# Run a test under sudo, it's equivalent to bazel build <test> && sudo
# bazel-bin/<test>. This works more reliably, because sometime bazel-bin path
# might not be accessible.
#
# This test script hangs if the test starts multiple threads
# (for example: go routines) and not stopped in time; but wont happen if the
# test is run under shell directly.

# Invoke sudo, otherwise bazel run will not handle password entry correctly.
sudo echo

# Disabling the "test sharding strategy" is needed, despite that we are using bazel run here.
bazel run --run_under=sudo --test_sharding_strategy=disabled "$@"

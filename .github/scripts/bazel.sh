#!/bin/bash -ex

# A wrapper to accept a space-separated string and turn it into array for bazel
# commands.
bazel "$@"

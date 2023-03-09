#!/bin/bash -ex

# A wrapper to accept a space-separated string and turn it into array for bazel
# commands.
# NOTE: --config=github-actions might not be compatible with certain bazel
# commands. Keep this note for future reference.
bazel "$@" --config=github-actions

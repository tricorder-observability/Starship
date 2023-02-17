#!/bin/bash -ex

query_string="kind(rule, rdeps(//src/..., set($@)))"
query_string="${query_string} except attr('tags', 'manual', //src/...)"

# Query affected rules under //src and print a string that can be write
# into Github environment variables.
#
# NOTE: Needs to change multiple lines to space-separated one-line string.
rules=$(bazel query --keep_going ${query_string} | tr '\n' ' ')
if [[ "${rules}" != "" ]]; then
  echo "all=${rules}"
  echo "any_affected=true"
fi

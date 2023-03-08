#!/bin/bash -x

# Need to remove the -e option.

# A wrapper of bazel test to accept a space-separated string and turn it into
# array for bazel commands.
#
# Also ignore the 4 return code, which means no test targets but test requested.
# We do not want to create complex bazel query, instead just ignore such
# failure.
bazel test --config=github-actions --flaky_test_attempts=3 --cache_test_results=no "$@"
exit_code="$?"
if [[ "${exit_code}" == "4" ]]; then
  exit 0
fi
exit ${exit_code}

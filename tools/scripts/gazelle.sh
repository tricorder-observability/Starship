#!/bin/bash -e

# Update go mode and use bazel gazelle to update all BUILD files.

# TODO(yzhao): Not clear why -compat=<go_version> is required, but otherwise
# `go mod tidy` would fail with
# github.com/tricorder/future/src/starship imports
#   golang.org/x/sync/errgroup loaded from golang.org/x/sync@v0.0.0-20190423024810-112230192c58,
#   but go 1.16 would select v0.0.0-20220722155255-886fb9371eb4
# It seems we need to pin go version.
go mod tidy -compat=1.18
bazel run //:gazelle -- update-repos -prune -from_file=go.mod \
  -to_macro=bazel/go_deps.bzl%go_deps
bazel run //:gazelle

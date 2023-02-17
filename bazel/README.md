# Bazel

This directory is not a bazel package, there is no BUILD.bazel file.
Files under this directory should be referenced as `//:bazel/<relative_path>`,
i.e. use the relative path under the root of this entire workspace.
Files in this directory can only be referenced in $ToT/WORKSPACE.

* `go_deps.bzl`: Are updated by gazelle, **DO NOT EDIT**

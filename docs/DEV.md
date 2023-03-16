# DEV

## Testing GitHub workflow changes

Change the trigger condition to be pull request on main branch, and test the change in a Pull Request.

## Update proto file's generated .pb.go source files

```
// Run this script, it will build all go_proto_library() targets,
// and copy the generated .pb.go to be along side of the .proto files.
tools/scripts/build_and_update_go_protos.sh
```

## Bazel & bazelisk

Starship uses [bazel](https://bazel.build/), an opinionated monorepo build
system open sourced by Google.

You should first read an overview of the
[major Bazel concepts](https://bazel.build/concepts/build-ref)
to get yourself familiarized.

Specifically, Starship uses [bazelisk](https://github.com/bazelbuild/bazelisk)
instead of bazel directly. Bazelisk reads `.bazeliskrc` to get the specified
bazel version, download and execute the designated version. This ensures bazel
version consistent with repo commit, see [bazelisk config](
https://github.com/bazelbuild/bazelisk#how-does-bazelisk-know-which-bazel-version-to-run)
for more details.

## Dotfiles

Various config files for git, and other tools are here (Top of Tree or ToT).

* `.bazelignore` config for bazel to ignore certain directory, see [bazelrc](https://bazel.build/run/bazelrc)
* `.commitlintrc.yml` commitlint config file, see [commitlint](https://github.com/conventional-changelog/commitlint)
* `.bazel_fix_commands.json` ibazel auto fix config, see [bazel-watcher output-runner](https://github.com/bazelbuild/bazel-watcher#output-runner)
* `.golangci.yml` golangci-lint config, see [golangci-lint](https://github.com/golangci/golangci-lint)

## Reverting breaking commit first then fix

Do not try to fix non-trivial bug introduced by a previous commit.
Revert it, and then fix in a measured pace with sufficient guardrail of normal development process.
Otherwise, you might be too stressful and make a lot of mistakes.

## Docker run --etnrypoint

```
docker run --rm -it --name agent --entrypoint bash <image>
```

## bazel -c dbg

`bazel -c dbg` adds `gcflags=-N -l` to go code, which disable optimization `-N`
and inlining `-l`.

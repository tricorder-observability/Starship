# Starship

Starship is a next-generation Observability platform built on eBPF+WASM.
Starship is to modern Observability on Kubernetes platform, as ChatGPT is
to consumer knownledge discovery.

eBPF enables instrumentation-free data collection, and WASM complements eBPF's
inability to perform complex data processing.

Starship is developed by [Tricorder Observability](https://tricorder.dev/),
proudly supported by [MiraclePlus](https://www.miracleplus.com/) and the Open Source
community.

[![Bazel build and test](https://github.com/tricorder-observability/starship/actions/workflows/build-and-test.yml/badge.svg?event=pull_request)](https://github.com/tricorder-observability/starship/actions/workflows/build-and-test.yml)
[![Lint Code Base](https://github.com/tricorder-observability/starship/actions/workflows/super-linter.yaml/badge.svg?event=pull_request)](https://github.com/tricorder-observability/starship/actions/workflows/super-linter.yaml)
[![Additional lint](https://github.com/tricorder-observability/starship/actions/workflows/additional_lint.yml/badge.svg?event=pull_request)](https://github.com/tricorder-observability/starship/actions/workflows/additional_lint.yml)

[![Twitter URL](https://img.shields.io/twitter/url/https/twitter.com/bukotsunikki.svg?style=plastic&label=Follow%20%40tricorder_o11y)](https://twitter.com/tricorder_o11y)
[![Slack Badge](https://img.shields.io/badge/Slack-4A154B?logo=slack&logoColor=fff&style=plastic&label=Join%20Tricorder)](https://join.slack.com/t/tricorderobse-mfl6648/shared_invite/zt-1oxqtq793-rRA03FN1YuyCiQrN_TrZoQ)

Starship is an eBPF-based observability platform
for Kubernetes.

Starship provides all things you'll need to
get started with Zero-Cost (or Zero-Friction) Observability,
where no effort other than [installation](https://tricorder-observability.github.io/helm-charts/)
is required to get the value of Observability.

Starship provides `Service Map`, the most valuable information
for understanding Cloud Native applications,
and numerous other data, analytic, and
visualization capabilities to satisfy
the full spectrum of your needs in running
and managing Cloud Native applications
on Kubernetes.

Starship comprises 4 components:
1. A data collection agent running as daemonset
2. A database for storing observability data
3. A visualization component for Observability data,
   we use Grafana
4. An API Server to manage all the above compnents

The core of starship is the tricorder agent,
which runs data collection modules written in
your favorite language, and are executed in eBPF+WASM.

You can write your own modules in
C/C++ (Go, Rust, and more languages are coming).

Tricorder agent supports all major frontend languages
of writing eBPF programs, including:
* [BCC](https://github.com/iovisor/bcc)
* [BPFtrace](https://github.com/iovisor/bpftrace)
* Rust ([readbpf](https://github.com/foniod/redbpf) [aya](https://github.com/aya-rs/aya))

Additionally, [libbpf](https://github.com/libbpf/libbpf)-style eBPF binary object files
are supported as well.

Tricorder agent also supports writing WASM in C/C++, Go, and Rust as well.

More details will be added for how to combine eBPF and WASM
together to build a complete data collection module.

## Components

TODO: Add components descriptions

## Building Starship

The development environment is based on Ubuntu.
The easiest way to get started with building Starship is to use the dev image:

```
git clone git@github.com/tricorder-observability/startship
cd starship
devops/dev_image/run.sh
# Inside the container
bazel build src/...
```

`devops/dev_image/run.sh` is a script that mounts the `pwd` (which is the root
of the cloned Starship repo) to `/starship` inside the dev image.

### Provision development environment on localhost
You can use Ansible to provision development environment on your localhost.
First install `ansible`:

```
sudo apt-get install ansible-core -y
git clone git@github.com/tricorder-observability/startship
cd starship
sudo devops/dev_image/ansible-playbook.sh devops/dev_image/dev.yaml
```

This installs a list of apt packages, and downloads and installs a list of other
tools from online.

Afterwards, you need source the env var file to pick up the PATH environment
variable (or put this into your shell's rc file):
```
source devops/dev_image/env.inc
```

Afterwards, run `bazel build src/...` to build all targets in the Starship repo.

### Bazel & bazelisk

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

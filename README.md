# Starship ![Starship stars](https://img.shields.io/github/stars/tricorder-observability/starship?style=social)

⚠️ **Starship is in active development; it's unstable, and should only be used for evaluation**  ⚠️

![image](https://user-images.githubusercontent.com/112656580/219543149-2e2bbebc-1891-4dcb-ba66-0f8b7f1bcd68.png)
![image](https://user-images.githubusercontent.com/112656580/219542981-5a4e5fb1-0603-4c0b-91e2-c94c36a92c0b.png)

🖖 Starship 🪐 is a next-generation Observability platform built on 🐝 eBPF➕WASM ![image](https://user-images.githubusercontent.com/112656580/219543881-046af389-ca10-4dda-b79a-a60088a1220a.png)

🚀 Starship is to modern Observability, as ChatGPT is to consumer knowledge discovery.
🐝 eBPF enables instrumentation-free data collection, and
![image](https://user-images.githubusercontent.com/112656580/219543881-046af389-ca10-4dda-b79a-a60088a1220a.png)
WASM complements eBPF's inability to perform complex data processing.

Starship is developed by [Tricorder Observability](https://tricorder.dev/),
proudly supported by [MiraclePlus](https://www.miracleplus.com/) and the Open Source
community.

[![Bazel build and test](https://github.com/tricorder-observability/Starship/actions/workflows/build-and-test.yml/badge.svg?branch=main&event=push)](https://github.com/tricorder-observability/Starship/actions/workflows/build-and-test.yml)
[![Release](https://github.com/tricorder-observability/Starship/actions/workflows/release.yaml/badge.svg)](https://github.com/tricorder-observability/Starship/actions/workflows/release.yaml)

![GitHub issues](https://img.shields.io/github/issues/tricorder-observability/starship)
![GitHub pull requests](https://img.shields.io/github/issues-pr/tricorder-observability/starship)

![Twitter Follow](https://img.shields.io/twitter/follow/tricorder_o11y?style=social)
[![Slack Badge](https://img.shields.io/badge/Slack-4A154B?logo=slack&style=social&label=Join%20Tricorder)](https://join.slack.com/t/tricorderobse-mfl6648/shared_invite/zt-1oxqtq793-rRA03FN1YuyCiQrN_TrZoQ)

## Building Starship
The easiest way to get started with building Starship is to use the dev image:

```
git clone git@github.com:<fork>/Starship.git
cd Starship
# Luanch dev image container
devops/dev_image/run.sh
# Inside the container
bazel build src/...
```

`devops/dev_image/run.sh` mounts the `pwd` (which is the root of the cloned Starship repo)
to `/starship` inside the dev image.

## Get Started

☸️ [Helm-charts](https://tricorder-observability.github.io/Starship),
install Starship on your Kubernetes cluster with helm.

We recommend [Minikube](https://minikube.sigs.k8s.io/docs/start/)
[**v1.24.0**](https://github.com/kubernetes/minikube/releases/tag/v1.24.0).
Starship deployment is broken on Kubernetes version 1.25 and newer version because of incompatbility
of the bundled kube prometheus stack using Pod Security Policy, which was removed in Kubenetes 1.25.
See [issues/258](https://github.com/tricorder-observability/Starship/issues/258).

```
minikube version
minikube version: v1.24.0
commit: 76b94fb3c4e8ac5062daf70d60cf03ddcc0a741b

# First start the minikube cluster, and make sure have at least 8 cpus and
# 8196 MB memory.
minikube start --profile=${USER} --cpus=8 --memory=8192

# Create a namespace for installing Starship.
# Do not use a different namespace, as our documentation uses this namespace
# consistently, and you might run into unexpected issues with a different
# namespace.
kubectl create namespace tricorder
kubectl config set-context --current --namespace=tricorder

# Add Starship's helm-charts and install Starship
helm repo add tricorder-starship https://tricorder-observability.github.io/Starship
helm install my-starship tricorder-starship/starship
```

You should see the following pods running on your cluster.
![image](https://user-images.githubusercontent.com/112656580/220381364-65bebd35-bf6d-4780-981b-be94c5464607.png)

More details can be found at [helm-charts installation](
https://github.com/tricorder-observability/Starship/tree/main/helm-charts).

Then follow the [CLI build and install](
https://github.com/tricorder-observability/starship/blob/main/src/cli/README.md#build-and-install)
to install `starship-cli`.

Then expose the API Server http endpoint with `kubectl port-forward`:
```
# This allows starship-cli accessing API Server with
# --api-address=localhost:8081
kubectl port-forward service/api-server 8081:80 -n tricorder
```

> DO NOT use the Web UI, as it's not working right now
> [issue/#80](https://github.com/tricorder-observability/starship/issues/80).

Then make sure you are the root of the Starship repo, and create a pre-built module:
```
starship-cli --api-address localhost:8081 module create \
    --bcc-file-path=modules/ddos_event/ddos_event.bcc \
    --wasm-file-path=modules/ddos_event/write_events_to_output.wasm \
    --module-json-path=modules/ddos_event/module.json
```
![image](https://user-images.githubusercontent.com/112656580/220375093-687b65b4-08fb-4be7-952a-89134306bb9c.png)

Then deploy this module:
```
starship-cli --api-address=localhost:8081 module deploy -i 0aa9e5db_ffce_4276_b37e_0b2dd82814a1
```
![image](https://user-images.githubusercontent.com/112656580/220375739-82f7b971-f0af-45e1-815e-e3c65c48be57.png)

```
starship-cli --api-address=localhost:8081 module deploy -i 0aa9e5db_ffce_4276_b37e_0b2dd82814a1
kubectl port-forward service/my-starship-grafana 8082:80 -n tricorder
```
Then open `http://localhost:8082`, login Grafana with username `admin` and password `tricorder`.
Then click the `Dashboards`->`Browse`, and then select the dashboard named `tricorder_<module_id>`.
You should see data reporting packets arriving with timestamp, as shown in the screenshot below.

![image](https://user-images.githubusercontent.com/112656580/220397224-5238110f-a1a0-4e0a-91de-4b9f9611caf9.png)

> Not yet very useful. We are working tirelessly 👩‍👨‍💻💻 on micro-service tracing!
> Stay tuned! 🫶

## Architecture

🤿 Before diving into the code base:

- Starship is built for Kubernetes platform. Starship provides all things you'll
  need to get started with Zero-Cost (or Zero-Friction) Observability.
- Starship provides `Service Map`, the most valuable information for
  understanding Cloud Native applications, and numerous other data, analytic,
  and visualization capabilities to satisfy the full spectrum of your needs in
  running and managing Cloud Native applications on Kubernetes.
- The core of starship is the tricorder agent, which runs data collection
  modules written in your favorite language, and are executed in eBPF+WASM.  You
  can write your own modules in C/C++ (Go, Rust, and more languages are coming).

We are working on supporting all major frontend languages of writing eBPF
programs, including:
* [BCC](https://github.com/iovisor/bcc)
* [BPFtrace](https://github.com/iovisor/bpftrace)
* Rust ([readbpf](https://github.com/foniod/redbpf)
  [aya](https://github.com/aya-rs/aya))

Additionally, [libbpf](https://github.com/libbpf/libbpf)-style eBPF binary
object files are supported as well.

## Components

* Starship [Tricorder](src/agent) (aka. Starship Agent): a data collection agent
  running as daemonset. Agent executes eBPF+WASM modules and export structured
  data to storage engine.  The code lives in [src/agent](src/agent).
* Starship [API Server](src/api-server): manages Tricorder agents, and Promscale
  & Grafana backend server; also supports management Web UI and CLI.  The code
  lives in [src/api-server](src/api-server).
* Starship [CLI](src/cli): the command line tool to use Starship on your
  Kubernetes cluster. The code lives in [src/cli](src/cli).
* Starship [Web UI](ui): a Web UI for using Starship.  The code lives in
  [ui](ui).

### 3rd party dependencies

* [Promscale](https://github.com/timescale/promscale): A unified metric and
  trace observability backend for Prometheus & OpenTelemetry.  Starship use
  `Promscale` to support Prom and OTel.
* [Grafana](https://github.com/grafana/grafana): Starship use `Grafana` to
  visualize Observability data.

### Prepherials

* [Kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) (KSM):
  listens to the Kubernetes API server and generates metrics about the state of
  the objects. Starship use `KSM` to expose cluster-level metrics.
* [Prometheus](https://github.com/prometheus/prometheus): collects metrics from
  `KSM` and then remote write to `Promscale`.
* [OpenTelemetry](https://github.com/open-telemetry): for distributed tracing
  and other awesome Observability features.

## Contributing

### Overview
- Fork the repo
- Createing Pull Request
- Ask for review

### Provision development environment on localhost
You can use Ansible to provision development environment on your localhost.
First install `ansible`:

```
sudo apt-get install ansible-core -y
git clone git@github.com:tricorder-observability/starship.git
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

### Creating Pull Requests

After making changes, run `tools/cleanup.sh` to cleanup the codebase, and then push
the changes to the forked repo, and create `Pull Request` on github Web UI.

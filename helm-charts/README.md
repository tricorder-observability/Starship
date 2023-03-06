# [Starship](https://github.com/tricorder-observability/starship) Helm charts

[![Release Helm Charts](https://github.com/tricorder-observability/Starship/actions/workflows/release-chart.yaml/badge.svg)](https://github.com/tricorder-observability/Starship/actions/workflows/release-chart.yaml)

[中文](./docs/README_CN.md)

![image](https://user-images.githubusercontent.com/112656580/219543149-2e2bbebc-1891-4dcb-ba66-0f8b7f1bcd68.png)
![image](https://user-images.githubusercontent.com/112656580/219542981-5a4e5fb1-0603-4c0b-91e2-c94c36a92c0b.png)

Helm charts for deploying
[Starship](https://github.com/tricorder-observability/starship)，
[Tricorder Observability](https://tricorder.dev)'s
next-generation Observability platform.

🖖 Starship 🪐 is a next-generation Observability platform built on
🐝 eBPF➕WASM ![image](https://user-images.githubusercontent.com/112656580/219543881-046af389-ca10-4dda-b79a-a60088a1220a.png)

🚀 Starship is to modern Observability, as ChatGPT is to consumer knownledge discovery.
🐝 eBPF enables instrumentation-free data collection, and
![image](https://user-images.githubusercontent.com/112656580/219543881-046af389-ca10-4dda-b79a-a60088a1220a.png)
WASM complements eBPF's inability to perform complex data processing.

Starship currently only runs on Kubernetes. Starship provides eBPF-powered
instrumentation-free Service Map.  No need to change a single line
of code in your application, instantly access a single-pane view of the
high-level status of your Cloud Native applications on Kubernetes.

Starship also collects data from [Prometheus](https://prometheus.io/) and
[OpenTelemetry](https://opentelemetry.io/).

Starship uses [Grafana](https://github.com/grafana/grafana) for visualization.
Use the following info to logon Grafana:
```text
username: admin
password: tricorder
```

## Prerequisites

TODO: Add instructions for other public Clouds.

### AWS EKS

- If you are using AWS EKS, install
  [EBS CSI](https://docs.aws.amazon.com/eks/latest/userguide/ebs-csi.html)
  on your EKS cluster. This is required because Helm charts create
  [PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
  for database pods, which requires EBS CSI.

### ALIYUN ACK

- If you are using Aliyun ACK, default install
  [Aliyun CSI](https://help.aliyun.com/document_detail/134722.html)
  on your ACK cluster. This is required because Helm charts create
  [PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
  for database pods, which requires Aliyun CSI.
- Before installing Starship, you need to check default storageclass,
  ```shell
  kubectl get storageclass | grep default
  ````
- If your cluster has no default storageclass,
you'll need to run the command below to create the default storageclass
  ```shell
  kubectl patch storageclass <you-storageclass-name> --patch '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
  # example:
  kubectl patch storageclass alibabacloud-cnfs-nas --patch '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
  ```
**WARNING:** Starship installation will fail if there is
no default storageclass on your cluster.

## Install

Change namespace to your own, here we use `tricorder` as an example.

```shell
helm repo add tricorder-stable \
    https://tricorder-observability.github.io/helm-charts
helm repo update
kubectl create namespace tricorder
helm install my-starship tricorder-stable/starship -n tricorder
```

## Access Starship web UI through LoadBalancer external IP

Starship by default expose web UI service through `LoadBalancer` service.
If your cluster has configured LoadBalancer that supports external access,
like AWS LoadBalancer Controller or an Ingress Controller, you can access the service
directly through `api-server` service's `ExteranIP`:

```shell
kubectl get service -n tricorder
```

![image](https://user-images.githubusercontent.com/112656580/215043391-6c4cd4bd-3a58-472f-a688-b88f11ef90c1.png)

Navigate to `http://${EXTERNAL-IP}` in your browser to access Starship's Web UI,
note that the protocol is **HTTP**, not **HTTPS**.
You will be able to open Grafana instance
by following the link on the left panel of the management Web UI.

## Access Starship Web UI through `kubectl port-forward`

You can expose the services to your local network
using the `kubectl port-forward` command.
Use the following command to expose the Starship managenment UI service
on port 18080 on localhost: `http://localhost:18080`.

```shell
kubectl -n tricorder port-forward service/api-server 18080:80
```

## Access Starship Web UI through Ingress

Ingress exposes HTTP and HTTPS routes from outside the cluster
to services with in the cluster.
Use the following command to expose the Starship managenment UI service
as `tstarship.io` host on ingress port 80.

```shell
kubectl apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: tricorder
  namespace: tricorder
spec:
  rules:
  - host: starship.io
    http:
      paths:
      - pathType: Prefix
        path: "/"
        backend:
          service:
            name: api-server
            port:
              number: 80
EOF

```
If you want to use more ingress features, such as TLS termination,
please follow this [documentation](https://kubernetes.io/docs/concepts/services-networking/ingress/)
to write ingress configuration.

## Caveats

**WARNING** Do not install multiple Starship releases in multiple namespaces.
That won't work because of system limitations. If you accidentally did that,
follow the [uninstall](#uninstall) instructions to remove all artifacts and reinstall.

**WARNING** Do not install multiple releases in the same namespace.

**WARNING:** This project is currently in active development. Consider this a
technical preview only.

## Data Retention

Metric and Trace data has an automated retention that drops data after a certain
age. The default retention is 7 days:

```yaml
promscale:
  config:
    startup.dataset.config: |
      metrics:
        compress_data: true
        default_retention_period: 7d
      traces:
        default_retention_period: 7d
```

and above retention can be customized by `--values` flag, We can change
`default_retention_period`'s value from `7 days` to `30 days`:

- create patch yaml file for custom values:

```shell
cat > rentention_patch.yaml << EOF
promscale:
  config:
    startup.dataset.config: |
      metrics:
        compress_data: true
        default_retention_period: 30d
      traces:
        default_retention_period: 30d
EOF
```

- override default settings by `--values` flag:

```shell
helm install my-starship tricorder-stable/starship -n tricorder \
    --values rentention_patch.yaml
```

## Uninstall

To uninstall a release you can run:

```shell
helm uninstall my-starship -n tricorder
```

After uninstalling helm release some objects will be left over. To remove them
follow next sections.

### Cleanup secrets

Secret's created with the deployment aren't deleted. These secrets need to be
manually deleted:

```shell
kubectl delete -n tricorder \
    $(kubectl get secrets -n tricorder -l "app=timescaledb" -o name)
```

### Cleanup configmap

```shell
kubectl delete -n tricorder \
    $(kubectl get configmap -n tricorder -l "app=my-starship-promscale" -o name)
```

### Cleanup Kube-Prometheus secret

One of the Kube-Prometheus secrets created with the deployment isn't deleted.
This secret needs to be manually deleted:

```shell
kubectl delete secret -n tricorder my-starship-kube-prometheus-stack-admission
```

### Cleanup DB PVCs and Backup

Removing the deployment does not remove the Persistent Volume Claims (pvc)
belonging to the release. For a full cleanup run:

```shell
kubectl delete -n tricorder \
    $(kubectl get pvc -n tricorder -l release=my-starship -o name)
```

### Prometheus PVCs

Removing the deployment does not remove the Persistent Volume Claims (pvc) of
Prometheus belonging to the release. For a full cleanup run:

```shell
kubectl delete -n tricorder $(kubectl get pvc -n tricorder \
  -l operator.prometheus.io/name=my-starship-kube-prometheus-stack-prometheus \
  -o name)
```

### Prometheus CRDs, ValidatingWebhookConfiguration and MutatingWebhookConfiguration

```shell
kubectl delete crd alertmanagerconfigs.monitoring.coreos.com \
    alertmanagers.monitoring.coreos.com \
    probes.monitoring.coreos.com \
    prometheuses.monitoring.coreos.com \
    prometheusrules.monitoring.coreos.com \
    servicemonitors.monitoring.coreos.com \
    thanosrulers.monitoring.coreos.com \
    podmonitors.monitoring.coreos.com
```

```shell
kubectl delete MutatingWebhookConfiguration my-starship-kube-promethe-admission
kubectl delete ValidatingWebhookConfiguration my-starship-kube-prometheus-admission
```

### Delete Namespace

```shell
kubectl delete namespace tricorder
```

## Advanced topics

### Send OpenTelemetry data to Starship

[Send OpenTelemetry data to Starship](./docs/STARSHIP_OTLP.md).

### Override default values

You can override configuration values defined in `charts/starship/Values.yaml`
with `--set` flags.

```shell
# Override Starship service's type to ClusterIP
helm upgrade my-starship tricorder-stable/starship -n tricorder \
    --set service.type=ClusterIP

# Override Starship container images' tag
helm upgrade my-starship tricorder-stable/starship -n tricorder \
    --set images.tag=<a specific tag>

# Override Starship container images' registry
helm upgrade my-starship tricorder-stable/starship -n tricorder \
    --set images.registry=<a specific imageRegistry>

# Override Starship api-server container image's tag
helm upgrade my-starship tricorder-stable/starship -n tricorder \
    --set apiServer.image.tag=<a specific tag>
```

### Install from local helm-charts repo

`git clone` this repo, and replace `tricorder-stable/starship` with the
charts local path `charts/starship`.

```shell
git clone git@github.com:tricorder-observability/helm-charts.git
cd helm-charts
# You'll need this step to fetch the dependent charts
helm dep update charts/starship
helm install my-starship charts/starship -n tricorder
```

All commands listed in the previous [install](#install) section works when you
swap `tricorder-stable/starship` with `charts/starship` when the PWD is the
root of the repo.

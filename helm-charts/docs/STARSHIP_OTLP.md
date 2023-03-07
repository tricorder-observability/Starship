# Send OpenTelemetry data to Starship

You can send traces to Starship using OTLP with any of the OpenTelemetry client
SDKs, instrumentation libraries, or the OpenTelemetry Collector.

You can quickly preview the effect by installing the
[OpenTelemetry demo](https://github.com/open-telemetry/opentelemetry-demo)
through [OpenTelemetry Demo Helm Chart](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-demo)

## Installing the OpenTelemetry Demo Chart

Before the demo installation, we need to make sure our starship components and
relative `kubernetes service` are ready, you check by following command:

e.g. To check our starship promescale service in `tricorder` namespace:

```shell
kubectl get service -n tricorder | grep promscale
```

So, `my-tricorder-promscale.tricorder.svc.cluster.local:9202` is the endpoint of
Promescale.

```shell
export STARSHIP_PROMESCALE_ENDPOINT=\
    "my-tricorder-promscale.tricorder.svc.cluster.local:9202"
```

Create patch yaml file for custom opentelemetry collector config:

```shell
cat > otel_demo_col_patch.yaml << EOF
opentelemetry-collector:
  config:
    receivers:
      otlp:
        protocols:
          grpc:
          http:
            cors:
              allowed_origins:
                - "http://*"
                - "https://*"

    exporters:
      otlp:
        endpoint: $STARSHIP_PROMESCALE_ENDPOINT
        tls:
          insecure: true
      prometheus:
        endpoint: '0.0.0.0:9464'

    processors:
      spanmetrics:
        metrics_exporter: prometheus

    service:
      pipelines:
        traces:
          processors: [memory_limiter, spanmetrics, batch]
          exporters: [otlp, logging]
        metrics:
          exporters: [prometheus, logging]
EOF
```

```shell
helm repo add open-telemetry \
    https://open-telemetry.github.io/opentelemetry-helm-charts
```

Create an kubernetes namespace for OpenTelemetry Demo installation:

```shell
kubectl create ns otel-demo
```

To install the chart with the release name `my-otel-demo`, run the following
command:

Since Starship Helm Chart already has Prometheus/Grafana installed and uses
TimeScaleDB to store trace data, you can disable duplicate components in
installation:

```shell
helm install my-otel-demo open-telemetry/opentelemetry-demo -n otel-demo \
    --set observability.jaeger.enabled=false \
    --set observability.prometheus.enabled=false \
    --set observability.grafana.enabled=false \
    --values otel_demo_col_patch.yaml
```

Last, check your opentelemetry demo pods, you should see a list of pods listed
below:
```shell
kubectl -n webstore-demo get pods
```
Then get back to the Starship web UI to examine the APM dashboards.
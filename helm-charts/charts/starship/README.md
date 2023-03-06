# Starship

Dependent external charts

* `timescaledb-single` Single node deployment of timescaleDB.
* `promscale` Timescale's prom & OTel connector, to enable ingesting and
  querying Prom & OTel data with the corresponding
  ingestion, transport, and query protocols.
* `kube-prometheus-stack` Prom collector, which defines how Prom collector
  connects to PromScale. This defines Prom collector.
  Application prometheus service monitor and pod monitor, this requires
  installing [Prom operator](
  https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/user-guides/getting-started.md#installing-the-operator).

## Templates

### Agent

* `daemonset.yaml` Defines agent daemonset, which runs on each and every
  Kubernetes node. They run eBPF+WASM data collection modules.  And collect
  Kubernetes process and container information, and writes to API Server for
  serialization.

TODO: Consider include [PGBouncer](https://github.com/pgbouncer/pgbouncer) to
have connection pool to support large kubernetes cluster.  Because agents right
now connects to PG directly for writing data collected by the eBPF+WASM data
collection modules.

### API Server with UI

* `service.yaml` Defines the service for connecting to the API server's backend
  server, and the management UI (backed by Nginx reverse proxy)
* `serviceaccount.yaml` Defines service accounts for API Sever's metadata
  service sub-component to access Kubernetes API server's updates of Kubernetes
  objects.
* `statefulset.yaml` Defines pods for API Server and management Web UI (nginx
  reverse proxy)

### Kubenetes Prometheus Stack

This directory includes the configurations of databases for storing prom, otel,
generic time-series data.

TODO: Rename to datasource or something relevant

* `connection-secret-job.yaml` Store credentials into K8s secret, for PromScale
  to connect to TimescaleDB.
* `grafana-dashboards-conf.yaml` Stores configurations of the pre-built
  dashboards.
* `grafana-datasources-sec.yaml` Stores configurations of prom, otel, and
  timescale time-series database query endpoints and credentials.

### Tricorder Database Initialization

* `post-init-configmap.yaml` Create Kuternetes ConfigMap that stores the content
  of scripts for initializing Timescale DB (aka Postgres + Timescale
  extensions).
* `timescaledb-extensions.yaml` Create TimescaleDB extensions for Porm and OTel

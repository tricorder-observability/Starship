# API Server

API Server manages Tricorder Agents, PromScale database, and Grafana backend.

It also supports the management Web UI for users to manage the data collection modules,
and instruct API server to deploy one or more modules.

In the future, API Server will connect with Cloud Manager, serving as the
interface of the logical component of individual K8s cluster, to export data
collected in the local cluster to the Cloud storage.


## Sqlite

Embeded in API Server to store eBPF+WASM modules' metadata and status information
during their lifetime.

TODO(zhihui): The docs below needs revision.

## Create tables
- code table
```sql
CREATE TABLE "code" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" TEXT,
  "ebpf" TEXT,
  "wasm" TEXT,
  "create_time" TEXT,
  "status" integer,
  "wasm_content" TEXT,
  "wasm_file_name" TEXT,
  "schema_name" TEXT,
  "schema_attr" TEXT,
  "panel" TEXT,
  "fn" TEXT
);
```
- datasource table, save create success grafana datasource
```sql
CREATE TABLE "datasource" (
  "id" INTEGER NOT NULL,
  "name" TEXT,
  "create_time" TEXT,
  "type" TEXT,
  "host" TEXT,
  "user" TEXT,
  "password" TEXT,
  "tables_name" TEXT,
  "uid" TEXT,
  "grafana_datasource_id" INTEGER,
  "database" TEXT,
  PRIMARY KEY ("id")
);
```

- ebpf table
```sql
CREATE TABLE "ebpf" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "code" TEXT,
  "event_size" TEXT,
  "perf_buffers" TEXT,
  "create_time" TEXT,
  "entry" TEXT,
  "return" TEXT,
  "code_id" INTEGER
);
```
- grafana_api table, save grafana api token
```sql
CREATE TABLE "grafana_api" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" TEXT,
  "create_time" TEXT,
  "api_key" TEXT,
  "auth_value" TEXT
);
```
- grafana_dashboard table, save grafana dashboard
```sql
CREATE TABLE "grafana_dashboard" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" TEXT,
  "create_time" TEXT,
  "uid" TEXT,
  "grafana_dashboard_id" INTEGER,
  "version" TEXT
);
```
- schema table, save postgresql table schema
```sql
CREATE TABLE "schema" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" TEXT,
  "create_time" TEXT,
  "uid" TEXT,
  "grafana_dashboard_id" INTEGER,
  "version" TEXT
);
```

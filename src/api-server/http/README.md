# HTTP

The HTTP service component of the API server. It supports the management Web UI,
and Starship CLI. It forwards user requests to API server's other internal
components for actual processing. Implemented with
[Gin](https://github.com/gin-gonic/gin).

## SQLite

You can use SQLite CLI to examine the pre-generated tricorder.db file,
this file is packaged into the API Server image, and pre-loaded with a list of
modules.

- [SQLite Home](https://www.sqlite.org/index.html)
- [Download SQLite](https://www.sqlite.org/download.html)

### Open database file or create databasefile
if database file(example test.db) not exist, will auto create test.db
```bash
sqlite3 test.db
```

### show tables
```sql
.table
```
or
```sql
sqlite3 test.db .table
```

### create tables
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

### save data
```sql
INSERT INTO "ebpf" ( "id", "code", "event_size", "perf_buffers", "create_time", "entry", "return", "code_id" )
VALUES
	( 10, 'b', '', 'sa', '2022-12-14 20:32:03', 's', 'asd', 27 );
```

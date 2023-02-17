#!/bin/bash

# Create database for Connector
echo "Creating database for echo http demo"

TRICORDER_NAMESPACE=tricorder
export TRICORDER_NAMESPACE
TIMESCALEDB_POD_NAME="$(kubectl get pod -n $TRICORDER_NAMESPACE -o name -l role=master,app=timescaledb)"
export TIMESCALEDB_POD_NAME

# shellcheck disable=SC2046
kubectl exec -n $TRICORDER_NAMESPACE -i --tty "${TIMESCALEDB_POD_NAME}" \
  -- psql -U postgres -c "CREATE DATABASE tricorder WITH OWNER postgres"
kubectl exec -n $TRICORDER_NAMESPACE -i --tty "${TIMESCALEDB_POD_NAME}" \
  -- psql -U postgres -c "\l"
kubectl exec -n $TRICORDER_NAMESPACE -i --tty "${TIMESCALEDB_POD_NAME}" \
  -- psql -U postgres -d tricorder -c \
  "CREATE TABLE IF NOT EXISTS http (
     time TIMESTAMP WITH TIME ZONE NOT NULL,
     id TEXT NOT NULL,
     method TEXT,
     proto TEXT,
     url TEXT,
     header TEXT,
     body TEXT,
     PRIMARY KEY(time, id)
   );
   SELECT create_hypertable('http', 'time', chunk_time_interval => INTERVAL '1 hour');
   SELECT add_retention_policy('http', INTERVAL '1 days');"

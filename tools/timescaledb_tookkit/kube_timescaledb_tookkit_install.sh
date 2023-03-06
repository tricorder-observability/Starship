#!/bin/bash

# Create database for Connector
echo "Creating timescaledb_toolkit for APM/otel demo"

TRICORDER_NAMESPACE=tricorder
export TRICORDER_NAMESPACE
TIMESCALEDB_POD_NAME="$(kubectl get pod -n $TRICORDER_NAMESPACE -o name -l role=master,app=timescaledb)"
export TIMESCALEDB_POD_NAME

# shellcheck disable=SC2046
kubectl exec -n $TRICORDER_NAMESPACE -i --tty "${TIMESCALEDB_POD_NAME}" \
  -- /bin/sh /docker-entrypoint-initdb.d/010_install_timescaledb_toolkit.sh

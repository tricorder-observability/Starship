# Skaffold

Run skaffold in `skaffold-tricorder` namespace:
```
minikube start --profile=${USER}
skaffold run -f tools/skaffold/skaffold.yaml -n skaffold-tricorder
```

Starship's k8s deployment spec for development: `deployment.yaml`.

NOTE: Must be kept consistent with the helm-charts in:
https://github.com/tricorder-observability/helm-charts

## Access postgres

To get access to postgres instance with psql:
```
kubectl exec -it postgres-<suffix> sh
psql -U postgres

# List all databases in the postgres instance
psql-# \l

# Connect to database `tricorder`
psql-# \c tricorder

# Show data tables on `tricorder` database
psql-# \dt
# Show more details of the data tables
psql-# \dt+
```

## TODO

TODO(yzhao): [Making skaffold work with helm](https://github.com/tricorder-observability/helm-charts/issues/56)

# Minikube

Minikube related documentation.

## Avoid conflicts

When starting minikube cluster, apply a name to avoid conflicting with other
users on the shared dev host.
```
# To start a cluster
minikube start --profile ${USER}

# To stop a cluster
minikube stop --profile ${USER}

# To remove a cluster
minikube delete --profile ${${USER}
```

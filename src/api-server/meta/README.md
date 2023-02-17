# Meta

Meta is Tricorder's metadata service, which collects object updates from
Kubernetes API server.

Right now it reads the `~/.kube/config` to get the connection information of the
Kubernetes API server, and connect to it.

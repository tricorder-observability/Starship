# EKS

## Create EKS cluster

`create.sh` This creates a cluster with name `dev-cluster-${USER}`

## Delete EKS cluster

```
eksctl delete cluster <your cluster name>
```

## Creating EKS cluster on AWS management console

TODO: Need to translate this process to `eksctl` command

* Create cluster
* Create node group, naming node-group after the cluster ${cluster-name}-ng-0
* Enable ssh by selecting a key-pair. There is one step in creating node group

## Update `kubectl` config for EKS cluster

```
aws eks update-kubeconfig --name ${CLUSTER_NAME} --region ${REGION}
```

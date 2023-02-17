#!/bin/bash -ex

echo "Creating EKS cluster dev-cluster-${USER}"
echo "with config file dev/ops/aws/eks/cluster.yaml"

# Create an eks cluster with the running user in the cluster name.
ToT=$(git rev-parse --show-toplevel)
sed "s/{{USER}}/${USER}/" ${ToT}/devops/aws/eks/cluster.yaml |\
  eksctl create cluster -f -

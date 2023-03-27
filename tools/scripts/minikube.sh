#!/bin/bash -ex

echo "Starting Minikube cluster ..."
minikube start -p ${USER} --cpus=8 --memory=8196

repo_name="tricorder-starship"
echo "Adding Starship helm chart repo as ${repo_name} ..."
helm repo add ${repo_name} \
  https://tricorder-observability.github.io/Starship
helm repo update ${repo_name}

NS="tricorder"
kubectl create namespace ${NS}
helm install my-starship tricorder-starship/starship --debug

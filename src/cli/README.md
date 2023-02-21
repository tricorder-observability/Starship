# CLI

The CLI tool to work with Starship Observability platform.

```shell
git clone https://github.com/Tricorder Observability/starship.git
cd starship
bazel build -c opt //src/cli
cp ./bazel-bin/src/cli/cli_/cli ~/bin/starship-cli
export PATH=~/bin:$PATH
starship-cli -h
```

# Usage

```shell
starship-cli -h

starship-cli module -h

API_SERVER_ADDRESS="a536bf891ad354e6dbbde9cfa80733d5-1223172260.ap-southeast-1.elb.amazonaws.com:8080"
starship-cli module list --api-address ${API_SERVER_ADDRESS}

# create module
starship-cli module create --api-address ${API_SERVER_ADDRESS} \
    -b modules/sample_json/sample_json.bcc \
    -w modules/sample_json/copy_input_to_output.wasm \
    -m modules/sample_json/module.json

# deploy module
starship-cli module deploy --api-address ${API_SERVER_ADDRESS} \
    -i <module_id>
```

Note: the value of `--api-address` is the address of starship
apiserver.`localhost:8080` is defalt. If you need to connect to remote starship
api server which deployed in EKS/Kubernetes, following ways maybe helpful:

-  Access Starship Api Server through `kubectl port-forward`

You can expose the services to your local network using the `kubectl
port-forward` command. Use the following command to expose the Starship api
server service on port 8080 on localhost: http://localhost:8080.

```shell
$ kubectl  port-forward -n tricorder svc/api-server 8080:8080
```

- Access Starship Api Server through LoadBalancer external IP

Starship by default expose web UI and api server service through LoadBalancer
service. If your cluster has configured LoadBalancer that supports external
access, like AWS LoadBalancer Controller or an Ingress Controller, you can
access the service directly through api-server service's ExteranIP:

```shell
kubectl -n tricorder get service api-server --output \
    jsonpath='{.status.loadBalancer.ingress[0].hostname}'

a536bf891ad354e6dbbde9cfa80733d5-1223172260.ap-southeast-1.elb.amazonaws.com
```

# Develpment

```bash
$ cd starship/src/cli

# help message
$ starship-cli -h

$ starship-cli module -h
# list modules
$ go run main.go module list  --api-address a536bf891ad354e6dbbde9cfa80733d5-1223172260.ap-southeast-1.elb.amazonaws.com:8080
# create module
$ go run main.go module create -b ~/starship/modules/sample_json/sample_json.bcc -w ~/starship/src/agent/wasm/testdata/copy_input_to_output.wasm -m ~/starship/src/cli/cmd/module/testdata/module.json --api-address a536bf891ad354e6dbbde9cfa80733d5-1223172260.ap-southeast-1.elb.amazonaws.com:8080
# deploy module 
$ go run main.go module deploy  --api-address a536bf891ad354e6dbbde9cfa80733d5-1223172260.ap-southeast-1.elb.amazonaws.com:8080 -i xxxxxx
```

Note: the value of `--api-address` is the address of starship apiserver.`localhost:8080` is defalt. If you need to connect to remote starship api server which deployed in EKS/Kubernetes, following ways maybe helpful:

-  Access Starship Api Server through `kubectl port-forward`
You can expose the services to your local network using the `kubectl port-forward` command. Use the following command to expose the Starship api server service on port 8080 on localhost: http://localhost:8080.

```shell
$ kubectl  port-forward -n tricorder svc/api-server 8080:8080
```

- Access Starship Api Server through LoadBalancer external IP
Starship by default expose web UI and api server service through LoadBalancer service. If your cluster has configured LoadBalancer that supports external access, like AWS LoadBalancer Controller or an Ingress Controller, you can access the service directly through api-server service's ExteranIP:

```shell
$  kubectl -n tricorder get service api-server --output jsonpath='{.status.loadBalancer.ingress[0].hostname}'
a536bf891ad354e6dbbde9cfa80733d5-1223172260.ap-southeast-1.elb.amazonaws.com
```

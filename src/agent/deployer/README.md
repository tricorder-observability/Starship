# Deployer

Deployer in agent accepts DeployModuleReq, and attemps to deploy eBPF+WASM modules,
and reports back the deployment result back to the ModuleDeployer server.

See `src/api-server/pb/service.proto` for ModuleDeployer service's definition.


## build and test

```bash
bazel run --run_under=sudo --test_sharding_strategy=disabled //src/agent/deployer:deployer_test
```

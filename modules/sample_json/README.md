# Sample JSON

Periodically output a const JSON string in eBPF perf event probe, and then copy
this JSON string to WASM and copy to output.

```
# Generate *.wasm file
make
# Remove *.wasm files
make clean
```

# Deploy module with starship-cli

- create module

```shell
$ starship-cli module create -b ./sample_json.bcc -w ./copy_input_to_output.wasm -m ./module.json --api-address your-api-server-address:8080

{"data":null,"code":"200","message":"create success, module id:c7bd055a_f34a_428c_bb22_e20e7df7edd6"}
```

- deploy module

```shell
$ starship-cli module deploy --api-address your-api-server-address:8080 -i c7bd055a_f34a_428c_bb22_e20e7df7edd6

{"data":null,"code":"200","message":"prepare deploy, please wait a moment."}
```

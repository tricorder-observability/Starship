# DDos Event

eBPF module: detects the DDOS event then submits the event to perf buffer
WASM module: transforms DDOS event output to JSON object

# Usage
```
# create ddos-event module
$ starship-cli  --api-address a6797719780714b1db070610c294b49c-821915578.ap-southeast-1.elb.amazonaws.com:8080 module create -b modules/ddos_event/ddos_event.bcc -w modules/ddos_event/cjson.wasm -m modules/ddos_event/module.json

data: []
code: "200"
message: 'create success, module id: 53b540e1_d661_4e5f_92aa_42f8a7b28da5'

# deploy module 
$ starship-cli  --api-address a6797719780714b1db070610c294b49c-821915578.ap-southeast-1.elb.amazonaws.com:8080 module -i 53b540e1_d661_4e5f_92aa_42f8a7b28da5

```
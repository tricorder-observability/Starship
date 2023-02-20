# DDos Event

eBPF module: detects the DDOS event, and then submits the event to perf buffer
WASM module: transforms DDOS event output to JSON

# Usage
```
# create module
starship-cli  --api-address a6797719780714b1db070610c294b49c-821915578.ap-southeast-1.elb.amazonaws.com:8080 \
    module create \
    -m modules/ddos_event/module.json \
    -b modules/ddos_event/ddos_event.bcc \
    -w modules/ddos_event/cjson.wasm

# deploy module
starship-cli  --api-address a6797719780714b1db070610c294b49c-821915578.ap-southeast-1.elb.amazonaws.com:8080 \
    module deploy \
    -i 53b540e1_d661_4e5f_92aa_42f8a7b28da5
```
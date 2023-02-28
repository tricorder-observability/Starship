
# starship load

load local module to SQLite db file

## Build binary from source


```shell
git clone https://github.com/tricorder-observability/starship.git
cd starship
bazel build -c opt //src/load
mkdir -p ~/bin
export PATH="~/bin:$PATH"
ln -s $(pwd)/bazel-bin/src/load/load_/load ~/bin/starship-load
starship-load -h
```

## Usage

```shell
# Write ddos_event module to the tricorder.db in this repo
starship-load load \
    -b modules/ddos_event/ddos_event.bcc \
    -w modules/ddos_event/write_events_to_output.wasm \
    -m modules/ddos_event/module.json \
    -o src/api-server/cmd/src/api-server/http/
```

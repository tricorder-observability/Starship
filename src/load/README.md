
# starship load

load local module to SQLite db file

## Build binary from source


```shell
git clone https://github.com/tricorder-observability/starship.git

cd starship

bazel build -c opt //src/load

cp ./bazel-bin/src/load/load_/load /usr/local/bin/starship-load
chmod +x /usr/local/bin/starship-load
starship-load -h
```

## Usage

```shell
starship-load -h

starship-load load -h

Flags:
  -b, --bcc-file-path string      The file path of bcc code.
  -o, --output string       The file path of SQLite db.
  -h, --help                      help for load
  -m, --module-file-path string   The file path of module in json format.
  -w, --wasm-file-path string     The file path of wasm code.
```

## Demo

```shell
starship-load load \
    -b modules/ddos_event/ddos_event.bcc \
    -w modules/ddos_event/cjson.wasm \
    -m modules/ddos_event/ddos_json.json \
    -o src/api-server/cmd/src/api-server/http/
```

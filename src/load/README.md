
# starship load

load local module to sqlit db file

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
  -d, --db-file-path string       The file path of sqlit db.
  -h, --help                      help for load
  -m, --module-file-path string   The file path of module in json format.
  -w, --wasm-file-path string     The file path of wasm code.
```

## Demo
```shell
./bazel-bin/src/load/load_/load load -b ~/src/starship/modules/ddos_event/ddos_event.bcc -w ~/src/starship/modules/ddos_event/cjson.wasm -m ~/src/starship/modules/ddos_event/ddos_json.json -d ~/src/starship/src/api-server/cmd/src/api-server/http/
```

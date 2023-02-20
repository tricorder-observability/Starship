# starship cli
The CLI way to manage starship observe modules.

# Installation
- Getting the binary for linux

TODO(jian): we need to create a CD pipeline to release this CLI binary, so that users can install `starship-cli` in binary way. 

```shell
#Check the release page:
#https://github.com/Tricorder Observability/starship/releases

export STARSHIP_VERSION=`curl https://github.com/Tricorder Observability/starship-cli/releases/latest  -Ls -o /dev/null -w %{url_effective} | grep -oE "[^/]+$"`
curl -LO https://github.com/Tricorder Observability/starship-cli/releases/download/$STARSHIP_VERSION/starship-cli_${STARSHIP_VERSION}_linux_amd64.tar.gz
tar -xvf starship-cli_${STARSHIP_VERSION}_linux_amd64.tar.gz  -C /usr/local/bin/

starship-cli -h
```

- Build binary from source

```shell
git clone https://github.com/Tricorder Observability/starship.git

cd starship

bazel build -c opt //src/cli

cp ./bazel-bin/src/cli/cli_/cli /usr/local/bin/starship-cli
chmod +x /usr/local/bin/starship-cli
starship-cli -h
```

# Usage

```shell
starship-cli -h

starship-cli module -h
```

## Development

See [development](./DEVELOPMENT.md) docs.

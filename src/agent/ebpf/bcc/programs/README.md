# eBPF

eBPF code. These code are pre-built eBPF code stored into Tricorder's module
database.

# Usage
### how to build test_probe?
```shell
cp cp test_uprobe.go.gen test_uprobe.go
## disabling inlining
go build  -gcflags="-l" -o test_probe test_probe.go
```
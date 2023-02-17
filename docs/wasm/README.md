# WASM Pre-study

**Index**

[WASM](https://tricorder.feishu.cn/wiki/wikcnzljLZ0AeiSYd721wEOPOsc)
includes the investigation on WASM runtime.

## TODO: still missing parts (for end to end demo)

* `1`: OK
* `2`: Compile programs written in high-level languages C++ into wasm
* `2.1`: OK Compile Rust into WASM
* `3`: OK
* `4`: OK
* `5`: OK
* `6`: OK
* `6.5`: call wasm functions with events data
* `7`: write metrics data into relational DB

```
                 bpf code     data processing code
                              (in some languages)
                  |               |
                  |               |2
                  |1              |
                  |             wasm code
                  |               |
                  |               |3
                  \/              \/
               +---------------------+                                 +------------+       +------------+
               |                     |  6  spwan sandbox, call func    |            |  7    |            |
               |  TriCorder  Agent   |--------------->-----------------|wasm sandbox|-------| Metrics DB |
               |                     | 6.5 with events data from kernel|            |       |            |
 Userspace     +---------------------+                                 +------------+       +------------+
                  |               /\
~~~~~~~~~~~~~~~~~~|~~~~~~~~~~~~~~~||~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
                  |4              ||5
 Kernel space     |               || events
                  \/              //
             +------------+      //
             | BPF Engine |=====//
             +------------+ 
```

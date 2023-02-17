# Wasm

WASM runtime

TODO: https://www.youtube.com/watch?v=DdDF_UZO5IQ&list=PPSV
Shared memory for sharing data between userspace and wasm runtime
Consider use the same idea.
It was mentioned they use serialization and de-serialization at the beginning,
then switched to shared memory
Actually, it's not clear how they signal data passing to userspace and wasm
runtime
wasmtime adds shared memory
Global pool of shared memory, so the exited function can be removed
automatically.

TODO: Investigate socket-based communication

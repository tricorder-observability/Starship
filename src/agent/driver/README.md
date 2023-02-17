# Driver

Driver manages modules. A module is a shortname for eBPF+WASM data collection
module. It has 2 components:

- An eBPF module for collecting data from inside Linux Kernel.
- A WASM module for processing the data obtained from the eBPF module, and
  produces structured output

Driver manages the whole lifetime of the modules:

- The whole lifetime process from deployment through undeployment and deletion.
- Data processing from polling data from eBPF, push them into WASM runtime,
  get structured output data from WASM, and writing the structured data
  to Postgres.
- Many other minor works.

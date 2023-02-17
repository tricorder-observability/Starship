# Agent

Tricorder Agent does most of the heavy lifting work:

* Deploy eBPF + WASM data collection module, and manage their lifetime
* Collecting data from data collection module(s)
* Exporting data to intra-cluster PromScale database (PromScale database is just Postgres DB with Timescale's time series extension, https://www.timescale.com/promscale)

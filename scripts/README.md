# Utility Scripts

The scripts in this folder are utilities for the MGC SDKs. Most of them are written in
Python following PEP convention enforced by `flake8` and `black`.

One can run most of them with:

```shell
python3 <script> -h
```

## Scripts

### [sync_oapi.py](./sync_oapi.py):

Sync OpenAPI specs between internal implementation, which is a JSON generated from the
actual backend implementation (always updated) and the current external OpenAPI being
served.

Some transformations that need to be done on top of the current external spec:

1. Fix `server.urls` to match the external URLs
2. Replace any `requestBody` that mismatches from the internal spec, update with the
internal one
3. Change the error object since externals have a different Kong formatting for it

#### Running

For help:

```shell
python3 sync_oapi.py  --help
```

For running:

```shell
python3 sync_oapi.py <url-to-internal-spec> <path-ext-yaml-spec> -o <output-path-new-ext>
```

Where url to internal spec is something like: "https://vm-region.proxy.com/openapi.json"

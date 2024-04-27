# Creates a new replica for an instance asynchronously.

## Usage:
```bash
Usage:
  ./cli dbaas replicas create [flags]
```

## Product catalog:
- Flags:
- --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
- --exchange string               Exchange (default "dbaas-internal")
- --flavor-id uuid                Flavor Id
- -h, --help                          help for create
- --name string                   Name (max character count: 255) (required)
- --source-id uuid                Source Id (required)
- -v, --version                       version for create

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-ne1")
- --server-url uri         Manually specify the server to use

## Flags:
```bash

```

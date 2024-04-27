# Returns a list of database instances for a x-tenant-id.

## Usage:
```bash
Usage:
  ./cli dbaas instances list [flags]
```

## Product catalog:
- Flags:
- --control.expand enum
- Instance extra attributes or relations to show with the main query. When available, more than one value
- can be informed using commas. e.g: '--control.expand="replicas"' (must be "replicas")
- --control.limit integer     Limit (range: 1 - 25) (default 10)
- --control.offset integer   Offset (min: 0)
- --exchange string          Exchange (default "dbaas-internal")
- -h, --help                     help for list
- --status enum              InstanceStatusResponse: An enumeration. (one of "ACTIVE", "BACKING_UP", "CREATING", "DELETED", "DELETING", "ERROR", "ERROR_DELETING", "MAINTENANCE", "PENDING", "REBOOT", "RESIZING", "RESTORING", "STARTING", "STOPPED" or "STOPPING")
- -v, --version                  version for list

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-ne1")
- --server-url uri         Manually specify the server to use

## Flags:
```bash

```

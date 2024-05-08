# List all backups.

## Usage:
```bash
Usage:
  ./mgc dbaas backups list [flags]
```

## Product catalog:
- Flags:
- --control.limit integer     Limit (range: 1 - 50) (default 10)
- --control.offset integer   Offset (min: 0)
- --exchange string          Exchange (default "dbaas-internal")
- -h, --help                     help for list
- --mode enum                BackupMode: An enumeration. (one of "FULL" or "INCREMENTAL")
- --status enum              BackupStatusResponse: An enumeration. (one of "CREATED", "CREATING", "DELETED", "DELETING", "ERROR", "ERROR_DELETING" or "PENDING")
- --type enum                BackupType: An enumeration. (one of "AUTOMATED" or "ON_DEMAND")
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


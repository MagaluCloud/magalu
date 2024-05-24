# Update backup status by id on urp

## Usage:
```bash
Usage:
  ./mgc virtual-machine-xaas backups urp [external-backup-id] [flags]
```

## Product catalog:
- Flags:
- --error string                Error
- --external-backup-id string   External Backup Id (required)
- -h, --help                        help for urp
- --min-disk integer            Min Disk (required)
- --size integer                Size (required)
- --state enum                  BackupState (one of "available", "deleted", "error" or "new") (required)
- --status enum                 BackupStatus (one of "completed", "creating", "deleted", "deleting", "error" or "provisioning") (required)
- -v, --version                     version for urp

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri         Manually specify the server to use

## Flags:
```bash

```


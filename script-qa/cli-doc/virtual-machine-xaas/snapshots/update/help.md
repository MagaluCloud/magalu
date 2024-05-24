# Internal route for update status of a snapshot when receive a update from URP.

## Usage:
```bash
### Note
This route is used only for internal proposes.
```

## Product catalog:
- Usage:
- ./mgc virtual-machine-xaas snapshots update [urp-id] [flags]

## Other commands:
- Flags:
- --error string    Error
- -h, --help            help for update
- --size integer    Size
- --state enum      SnapshotStateV1 (one of "available", "deleted", "new" or "not_used")
- --status enum     SnapshotStatusV1 (one of "completed", "creating", "creating_error", "creating_error_quota", "deleted", "deleted_error", "deleting" or "provisioning")
- --urp-id string   Urp Id (required)
- -v, --version         version for update

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri         Manually specify the server to use
```


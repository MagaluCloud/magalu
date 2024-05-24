# Internal route for update status of a image when receive a update from URP.

## Usage:
```bash
### Note
This route is used only for internal proposes.
```

## Product catalog:
- Usage:
- ./mgc virtual-machine-xaas images urp [urp-id] [flags]

## Other commands:
- Flags:
- --error string    Error
- -h, --help            help for urp
- --status enum     ImageV1Status (one of "active", "creating", "deleted", "deprecated", "error", "importing" or "pending")
- --urp-id string   Urp Id (required)
- -v, --version         version for urp

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri         Manually specify the server to use
```


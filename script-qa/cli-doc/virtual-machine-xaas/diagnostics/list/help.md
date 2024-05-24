# Get a list of instances informing a tenant ID.

## Usage:
```bash
This route only list instances inserted on virtual machine database.
```

## Product catalog:
- Usage:
- ./mgc virtual-machine-xaas diagnostics list [flags]

## Other commands:
- Flags:
- --control.limit integer     Limit: limit the number of the results (max: 2147483647) (default 50)
- --control.offset integer    Offset: pagination for the results limited (range: 0 - 2147483647)
- --control.sort string       Sort: order of the results using informed fields (pattern: ^(^[\w-]+:(asc|desc)(,[\w-]+:(asc|desc))*)?$) (default "created_at:asc")
- -h, --help                     help for list
- --project-type enum        Project Type (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (default "default")
- -v, --version                  version for list

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri         Manually specify the server to use
```


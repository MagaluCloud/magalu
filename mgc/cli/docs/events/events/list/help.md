# Lists all cloud events emitted by other products.

## Usage:
```bash
Usage:
  ./mgc events events list [flags]
```

## Product catalog:
- Examples:
- ./mgc events events list --data='{"data.machine_type.name":"cloud-bs1.xsmall","data.tenant_id":"00000000-0000-0000-0000-000000000000"}'

## Other commands:
- Flags:
- --authid string            Authid
- --control.limit integer     Limit (max: 2147483647) (default 50)
- --control.offset integer    Offset (max: 2147483647)
- --data object              The raw data event
- Use --data=help for more details (default {})
- -h, --help                     help for list
- --id string                Id
- --product-like string      Product  Like
- --source-like string       Source  Like
- --tenantid string          Tenantid
- --time date-time           Time
- --type-like string         Type  Like
- -v, --version                  version for list

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri         Manually specify the server to use
      --x-tenant-id string     X-Tenant-Id
```


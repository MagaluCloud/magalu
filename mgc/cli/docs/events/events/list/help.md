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
      --api-key string           Use your API key to authenticate with the API
  -U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                                 use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                                 a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
  -t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                                 Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
      --debug                    Display detailed log information at the debug level
      --env enum                 Environment to use (one of "pre-prod" or "prod") (default "prod")
      --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
  -o, --output string            Change the ouput format. Use '--output=help' to know more details. (default "yaml")
  -r, --raw                      Output raw data, without any formatting or coloring
      --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri           Manually specify the server to use
      --x-tenant-id string       X-Tenant-Id
```


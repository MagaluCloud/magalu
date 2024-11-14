# List Virtual Machine instances

## Usage:
```bash
Usage:
  ./mgc virtual-machine instances list [flags]
```

## Product catalog:
- Flags:
- --control.limit integer     Limit: limit the number of the results (max: 1000) (default 50)
- --control.offset integer    Offset: pagination for the results limited (max: 2147483647)
- --control.sort string       Sort: order of the results using informed fields (pattern: ^(^[\w-]+:(asc|desc)(,[\w-]+:(asc|desc))*)?$) (default "created_at:asc")
- --expand array(string)     Expand: You can get more detailed info about: ['image', 'machine-type', 'machine-types', 'network', 'labels']  (default [])
- -h, --help                     help for list
- -v, --version                  version for list

## Other commands:
- Global Flags:
- --api-key string           Use your API key to authenticate with the API
- -U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
- use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
- a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
- -t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
- Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
- --debug                    Display detailed log information at the debug level
- --env enum                 Environment to use (one of "pre-prod" or "prod") (default "prod")
- --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
- -o, --output string            Change the output format. Use '--output=help' to know more details.
- -r, --raw                      Output raw data, without any formatting or coloring
- --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri           Manually specify the server to use

## Flags:
```bash

```


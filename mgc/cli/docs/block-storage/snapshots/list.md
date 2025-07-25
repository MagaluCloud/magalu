---
sidebar_position: 1
---
# List

Retrieve a list of Snapshots for the currently authenticated tenant.

## Usage:
```
mgc block-storage snapshots list [flags]
```

## Flags:
```
    --control.limit integer     Limit
    --control.offset integer    Offset
    --control.sort string       Sort (pattern: ^(^[\w-]+:(asc|desc)(,[\w-]+:(asc|desc))*)?$)
    --expand array(enum)       Expand
-h, --help                     help for list
    --name string              Name
    --type enum                Type (one of "instant" or "object")
```

## Global Flags:
```
    --api-key string           Use your API key to authenticate with the API
-U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                               use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                               a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
-t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                               Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
    --debug                    Display detailed log information at the debug level
    --env enum                 Environment to use (one of "pre-prod" or "prod") (default "prod")
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. Use '--output=help' to know more details.
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
```


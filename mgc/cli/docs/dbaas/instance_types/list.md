---
sidebar_position: 1
---
# List

Returns a list of available instance types. An instance type is a hardware template that defines the size of RAM and vCPU.

## Usage:
```
mgc dbaas instance-types list [flags]
```

## Flags:
```
    --compatible-product enum   Instance Type Compatible per product: Value referring to instance type compatible products. (one of "CLUSTER", "SINGLE_INSTANCE" or "SINGLE_INSTANCE_REPLICA")
    --control.limit integer     The maximum number of items per page. (range: 1 - 200)
    --control.offset integer    The number of items to skip before starting to collect the result set. (min: 0)
    --engine-id uuid            Value referring to engine Id.
-h, --help                      help for list
    --status enum               An enumeration. (one of "ACTIVE" or "DEPRECATED")
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


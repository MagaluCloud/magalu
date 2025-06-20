---
sidebar_position: 1
---
# List

Returns a list of database clusters for a x-tenant-id.

## Usage:
```
mgc dbaas clusters list [flags]
```

## Examples:
```
mgc dbaas clusters list --status="ACTIVE"
```

## Flags:
```
    --control.limit integer     The maximum number of items per page. (range: 1 - 25)
    --control.offset integer    The number of items to skip before starting to collect the result set. (min: 0)
    --engine-id uuid            Value referring to engine Id.
-h, --help                      help for list
    --parameter-group-id uuid   Value referring to parameter group Id.
    --status enum               Instance Status: Value referring to cluster status. (one of "ACTIVE", "BACKING_UP", "BALANCING", "CREATING", "DELETED", "DELETING", "ERROR", "ERROR_DELETING", "PENDING", "STARTING", "STOPPED" or "STOPPING")
    --volume.size integer       Volume.Size: Value referring to volume size.
    --volume.size-gt integer    Volume.Size Gt: Value referring to volume size greater than.
    --volume.size-gte integer   Volume.Size Gte: Value referring to volume size greater than or equal to.
    --volume.size-lt integer    Volume.Size Lt: Value referring to volume size less than.
    --volume.size-lte integer   Volume.Size Lte: Value referring to volume size less than or equal to.
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


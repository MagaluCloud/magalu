---
sidebar_position: 1
---
# List

List every value of a tag.

## Usage:
```
mgc tag values list [tag-name] [flags]
```

## Flags:
```
    --control.limit integer    Maximum number of items to return per page (range: 1 - 100)
    --control.offset integer   Number of items to skip before the first result (min: 0)
    --control.sort string      Sort criteria for the result
-h, --help                     help for list
    --name string              Name of the value to filter by (between 1 and 255 characters and pattern: ^[\w\ \-\[\]\(\)\.\:]+$)
    --tag-name string          Tag name that owns the value (required)
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
-o, --output string            Change the output format. You can use 'yaml', 'json' or 'table'.
-r, --raw                      Output raw data, without any formatting or coloring
    --server-url uri           Manually specify the server to use
```


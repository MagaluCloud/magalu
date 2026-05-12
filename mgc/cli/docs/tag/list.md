---
sidebar_position: 1
---
# List

List tags

## Usage:
```
mgc tag list [flags]
```

## Examples:
```
mgc tag list --color="ffffff00" --kinds='["finops"]'
```

## Flags:
```
    --color string             Color: 8-character lowercase hexadecimal string representing RGBA, without '#' prefix. (between 6 and 6 characters and pattern: ^[0-9a-fA-F]+$)
    --control.limit integer     Limit (range: 1 - 100)
    --control.offset integer    Offset (min: 0)
    --control.sort string       Sort
-h, --help                     help for list
    --kinds array(enum)        Kinds: a tag kind describe the purpose for this tag in mgc, the most common case is kind finops.
    --name string              name for tag (between 1 and 72 characters and pattern: ^[\w\ \-\[\]\(\)\.\:]+$)
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


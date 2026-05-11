---
sidebar_position: 2
---
# Create

Create tag, if values are informed, they will be created either

## Usage:
```
mgc tag create [flags]
```

## Examples:
```
mgc tag create --color="ffffff00" --kinds='["finops"]' --name="kubernetes-expenses"
```

## Flags:
```
    --color string           Color: 8-character lowercase hexadecimal string representing RGBA, without '#' prefix. (between 8 and 8 characters and pattern: ^[0-9a-f]+$)
    --force-authentication   Force authentication by sending the header even if this API doesn't require it
-h, --help                   help for create
    --kinds array(enum)      Kinds: a tag kind describe the purpose for this tag in mgc, the most common case is kind finops.
    --name string            name for tag (required)
    --values array(object)   Values (at most 20 items)
                             Use --values=help for more details
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


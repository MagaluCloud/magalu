---
sidebar_position: 5
---
# Update

Update the color, kinds or description of a tag.

## Usage:
```
mgc tag update [name] [flags]
```

## Examples:
```
mgc tag update --color="ffffff" --description="tag to monitor expenses with environments" --kinds='["finops"]' --name="kubernetes-expenses"
```

## Flags:
```
    --color string         Color of the tag, as a 6-digit hex RGB code without the '#' prefix (between 6 and 6 characters and pattern: ^[0-9a-fA-F]+$)
    --description string   Description of the tag (max character count: 500)
-h, --help                 help for update
    --kinds array(enum)    Kinds that describe what the tag is for, such as finops
    --name string          Name of the tag (required)
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


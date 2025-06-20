---
sidebar_position: 2
---
# Create

Select the scopes that the new API Key will have access to and set an expiration date

## Usage:
```
mgc auth api-key create [name] [description] [expiration] [flags]
```

## Examples:
```
mgc auth api-key create --expiration="2024-11-07 (YYYY-MM-DD)"
```

## Flags:
```
    --description string   Description of new api key
    --expiration string    Date to expire new api
-h, --help                 help for create
    --name string          Name of new api key (required)
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
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. Use '--output=help' to know more details.
-r, --raw                      Output raw data, without any formatting or coloring
```


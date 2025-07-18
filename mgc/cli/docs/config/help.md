---
sidebar_position: 0
---
# Config

Configuration values are available to be set so that they persist between
different executions of the MgcSDK. They reside in a YAML file when set.
Config values may also be loaded via Environment Variables. Any Config available
(see 'list') may be exported as an env variable in uppercase with the 'MGC_' prefix

## Usage:
```
mgc config [flags]
mgc config [command]
```

## Commands:
```
delete      Delete/unset a Config value that had been previously set
get         Get a specific Config value that has been previously set
get-schema  Get the JSON Schema for the specified Config
list        List all available Configs
set         Set a specific Config value in the configuration file
```

## Flags:
```
-h, --help   help for config
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


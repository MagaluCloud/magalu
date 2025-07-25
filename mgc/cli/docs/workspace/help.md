---
sidebar_position: 0
---
# Workspace

Workspace hold auth and runtime configuration, like tokens and log filter settings.
Users can create as many workspaces as they choose to. Auth and config operations will affect only the
current workspace, so users can alter and switch between workspaces without loosing the previous configuration

## Usage:
```
mgc workspace [flags]
mgc workspace [command]
```

## Commands:
```
create      Creates a new workspace
delete      Deletes the workspace with the specified name
get         Get current workspace.
list        List all available workspaces
set         Sets workspace to be used
```

## Flags:
```
-h, --help   help for workspace
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


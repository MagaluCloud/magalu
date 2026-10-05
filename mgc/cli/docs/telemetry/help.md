---
sidebar_position: 0
---
# Telemetry

Manage the pseudonymous usage data collected by the CLI.

After each command the CLI sends one event with the command path, the names of
the flags that were set (never their values), the outcome and a failure category,
the duration, the CLI version, the operating system, the install method, the
execution context (interactive, ci or agent), the tenant ID, the last request ID
and an anonymous installation ID. Flag values, credentials, error messages and
command output are never sent.

To opt out, run 'mgc telemetry disable' or set MGC_CLI_TELEMETRY_OPTOUT=1 or DO_NOT_TRACK=1.

Privacy policy: https://magalu.cloud/termos-legais/politica-de-privacidade/

## Usage:
```
mgc telemetry [command]
```

## Commands:
```
disable     Disable usage data collection
enable      Enable usage data collection
status      Show whether usage data collection is enabled
```

## Flags:
```
-h, --help   help for telemetry
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
-o, --output string            Change the output format. You can use 'yaml', 'json' or 'table'.
-r, --raw                      Output raw data, without any formatting or coloring
```


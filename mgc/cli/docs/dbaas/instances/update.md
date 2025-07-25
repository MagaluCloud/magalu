---
sidebar_position: 8
---
# Update

Updates a database instance.

## Usage:
```
mgc dbaas instances update [instance-id] [flags]
```

## Examples:
```
mgc dbaas instances update --backup-retention-days=7 --backup-start-at="04:00:00" --parameter-group-id="44ae8773-a21e-4d5e-a38f-b677ccfeb7f8"
```

## Flags:
```
    --backup-retention-days integer   Backup Retention Days: The number of days that a particular backup is kept until its deletion.
    --backup-start-at time            Backup Start At: Start time (UTC timezone) which is allowed to start the automated backup process.
    --cli.list-links enum[=table]     List all available links for this command (one of "json", "table" or "yaml")
    --deletion-protected              Deletion Protected
-h, --help                            help for update
    --instance-id uuid                Value referring to instance Id. (required)
    --parameter-group-id uuid         Parameter group Id
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


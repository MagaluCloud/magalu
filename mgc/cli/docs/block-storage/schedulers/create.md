---
sidebar_position: 2
---
# Create

Creates a schedule for snapshot creation.

## Usage:
```
mgc block-storage schedulers create [flags]
```

## Flags:
```
    --cli.list-links enum[=table]              List all available links for this command (one of "json", "table" or "yaml")
    --description string                       Description
-h, --help                                     help for create
    --name string                              Name (required)
    --policy object                            Policy (properties: frequency and retention_in_days)
                                               Use --policy=help for more details (required)
    --policy.frequency object                  Policy: Frequency (single property: daily)
                                               Use --policy.frequency=help for more details
                                               This is the same as '--policy=frequency:object'.
    --policy.frequency.daily object            Frequency: DailyFrequency (single property: start_time)
                                               Use --policy.frequency.daily=help for more details
                                               This is the same as '--policy.frequency=daily:object'.
    --policy.frequency.daily.start-time time   DailyFrequency: Start Time
                                               This is the same as '--policy.frequency.daily=start_time:time'.
    --policy.retention-in-days integer         Policy: Retention In Days (min: 1)
                                               This is the same as '--policy=retention_in_days:integer'.
    --snapshot object                          Snapshot (single property: type)
                                               Use --snapshot=help for more details (required)
    --snapshot.type enum                       Snapshot: SnapshotType (one of "instant" or "object")
                                               This is the same as '--snapshot=type:enum'.
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


---
sidebar_position: 2
---
# Create

Create a Snapshot for the currently authenticated tenant.

## Usage:
```
mgc block-storage snapshots create [flags]
```

## Examples:
```
mgc block-storage snapshots create --source-snapshot.id="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" --source-snapshot.name="some_resource_name" --volume.id="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" --volume.name="some_resource_name"
```

## Flags:
```
    --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
    --description string            Description (required)
-h, --help                          help for create
    --name string                   Name (between 3 and 50 characters) (required)
    --source-snapshot object        Source Snapshot (at least one of: single property: id or single property: name)
                                    Use --source-snapshot=help for more details
    --source-snapshot.id string     Source Snapshot: Id (min character count: 1)
                                    This is the same as '--source-snapshot=id:string'.
    --source-snapshot.name string   Source Snapshot: Name (between 1 and 255 characters)
                                    This is the same as '--source-snapshot=name:string'.
    --type enum                     SnapshotType (one of "instant" or "object")
    --volume object                 Volume (at least one of: single property: id or single property: name)
                                    Use --volume=help for more details
    --volume.id string              Volume: Id (min character count: 1)
                                    This is the same as '--volume=id:string'.
    --volume.name string            Volume: Name (between 1 and 255 characters)
                                    This is the same as '--volume=name:string'.
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


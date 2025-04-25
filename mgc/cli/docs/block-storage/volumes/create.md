# Create

Create a Volume for the currently authenticated tenant.

## Usage:
```
mgc block-storage volumes create [flags]
```

## Examples:
```
mgc block-storage volumes create --snapshot.id="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" --snapshot.name="some_resource_name" --type.id="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" --type.name="some_resource_name"
```

## Flags:
```
    --availability-zone string      Availability Zone
    --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
    --encrypted                     Indicates if the volume is encrypted. Default is False.
-h, --help                          help for create
    --name string                   Name (between 3 and 50 characters) (required)
    --size integer                  Size: Gibibytes (GiB) (range: 10 - 2147483648) (required)
    --snapshot object               Snapshot (at least one of: single property: id or single property: name)
                                    Use --snapshot=help for more details
    --snapshot.id string            Snapshot: Id (min character count: 1)
                                    This is the same as '--snapshot=id:string'.
    --snapshot.name string          Snapshot: Name (between 1 and 255 characters)
                                    This is the same as '--snapshot=name:string'.
    --type object                   Type (at least one of: single property: id or single property: name)
                                    Use --type=help for more details (required)
    --type.id string                Type: Id (min character count: 1)
                                    This is the same as '--type=id:string'.
    --type.name string              Type: Name (between 1 and 255 characters)
                                    This is the same as '--type=name:string'.
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


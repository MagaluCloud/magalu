# List

Returns a list of ports for a provided vpc_id and x-tenant-id. The list will be paginated, it means you can easily find what you need just setting the page number(_offset) and the quantity of items per page(_limit). The level of detail can also be set

## Usage:
```
mgc network vpcs ports list [vpc-id] [flags]
```

## Flags:
```
    --control.limit integer        Items Per Page (min: 1) (default 10)
    --control.offset integer       Page Number (min: 1) (default 1)
    --detailed                     Detailed (default true)
-h, --help                         help for list
    --name string                  Name of the port to list: Filter ports results with name
    --port-id-list array(string)   Port Id List (default [])
    --vpc-id string                vpc_id: ID of VPC to list ports (required)
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


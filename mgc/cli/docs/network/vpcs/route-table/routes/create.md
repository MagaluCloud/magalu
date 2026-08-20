---
sidebar_position: 2
---
# Create

Add a route to a VPC's route table.

## Usage:
```
mgc network vpcs route-table routes create [vpc-id] [flags]
```

## Examples:
```
mgc network vpcs route-table routes create --targets.id="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" --targets.type="vpc_peering"
```

## Flags:
```
    --cidr-destination string       Cidr Destination (required)
    --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
    --description string            Description
-h, --help                          help for create
    --targets object                TargetSchema (properties: id and type)
                                    Use --targets=help for more details (required)
    --targets.id uuid4              TargetSchema: Id
                                    This is the same as '--targets=id:uuid4'.
    --targets.type enum             TargetSchema: RouteTargetType (one of "port_id" or "vpc_peering")
                                    This is the same as '--targets=type:enum'.
    --vpc-id string                 Vpc Id: ID of the VPC whose route table receives this route (the source side of the traffic). (required)
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
    --iam-project-id string    Project scope for IAM commands only. Overrides the configured IAM project
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. You can use 'yaml', 'json' or 'table'.
    --project-id string        Project to scope the requests. Overrides the configured project
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
```


# Create

Create Backend

## Usage:
```
mgc load-balancer network-backends create [load-balancer-id] [flags]
```

## Flags:
```
    --balance-algorithm string      Balance Algorithm: The load balancing algorithm used by the backend (e.g., round_robin) (required)
    --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
    --description string            A brief description of the backend
    --health-check-id string        Health Check Id (at least one of: uuid)
-h, --help                          help for create
    --load-balancer-id uuid         load_balancer_id: ID of the attached Load Balancer (required)
    --name string                   The unique name of the backend (max character count: 64) (required)
    --targets array                 Targets: The list of target configurations for the backend (at least one of: array or array)
                                    Use --targets=help for more details (default [])
    --targets-type enum             Targets Type: The type of targets used by the backend (e.g., instance, raw) (one of "instance" or "raw") (required)
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


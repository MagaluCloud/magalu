---
sidebar_position: 2
---
# Create

Create HTTP Route Rule

## Usage:
```
mgc beta load-balancer application routes rules create [load-balancer-id] [route-id] [flags]
```

## Examples:
```
mgc beta load-balancer application routes rules create --headers='[{"name":"X-Env","value":"prod"}]'
```

## Flags:
```
    --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
    --headers array(object)         List of HTTP headers required for matching.
                                    Use --headers=help for more details
-h, --help                          help for create
    --load-balancer-id uuid         load_balancer_id: ID of the attached Load Balancer (required)
    --method enum                   HttpMethodEnum (one of "DELETE", "GET", "PATCH", "POST" or "PUT")
    --prefix string                 URI prefix to be routed (e.g., /api).
    --route-id uuid                 route_id: ID of the HTTP route (required)
```

## Global Flags:
```
    --api-key string           Use your API key to authenticate with the API
    --backend-id uuid          backend-id: ID of the backend this rule points to (required)
-U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                               use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                               a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
-t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                               Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
    --debug                    Display detailed log information at the debug level
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. You can use 'yaml', 'json' or 'table'.
    --project-type enum        ProjectType: Specifies the project type to which the load balancer belongs (one of "dbaas", "default" or "k8saas")
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
    --zone enum                Availability zone where the resource will be created (one of "br-se1-a", "br-se1-b" or "br-se1-c") (required)
```


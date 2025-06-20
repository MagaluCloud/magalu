---
sidebar_position: 5
---
# Replace

Update Load Balancer by ID

## Usage:
```
mgc load-balancer network-loadbalancers replace [load-balancer-id] [flags]
```

## Examples:
```
mgc load-balancer network-loadbalancers replace --backends='[{"health_check_id":"00000000-0000-0000-0000-000000000000","id":"00000000-0000-0000-0000-000000000000","targets":[{"id":"00000000-0000-0000-0000-000000000000","nic_id":"00000000-0000-0000-0000-000000000000","port":8080},{"id":"00000000-0000-0000-0000-000000000001","nic_id":"00000000-0000-0000-0000-000000000000","port":8080}],"targets_type":"instance"}]' --health-checks='[{"healthy_status_code":200,"healthy_threshold_count":8,"id":"00000000-0000-0000-0000-000000000000","initial_delay_seconds":30,"interval_seconds":30,"path":"/health-check","port":5000,"protocol":"tcp","timeout_seconds":10,"unhealthy_threshold_count":3}]' --tls-certificates='[{"certificate":"SGVsbG8sIFdvcmxkIQ==","id":"00000000-0000-0000-0000-000000000000","private_key":"SGVsbG8sIFdvcmxkIQ=="}]'
```

## Flags:
```
    --backends array(object)           Backends: The list of updated backend configurations
                                       Use --backends=help for more details
    --cli.list-links enum[=table]      List all available links for this command (one of "json", "table" or "yaml")
    --description string               The updated description of the load balancer (at least one of: string)
    --health-checks array(object)      Health Checks: The list of updated health check configurations
                                       Use --health-checks=help for more details
-h, --help                             help for replace
    --load-balancer-id uuid            load_balancer_id: ID of the Load Balancer to update (required)
    --name string                      The updated name of the load balancer (at least one of: string)
    --panic-threshold integer          Panic Threshold: Minimum percentage of failed upstreams that load balancer will consider to give an alert (range: 0 - 100)
    --tls-certificates array(object)   The list of updated TLS certificates
                                       Use --tls-certificates=help for more details
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


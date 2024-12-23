# Update Health Check by ID

## Usage:
```bash
Usage:
  ./mgc lbaas network-healthchecks replace [load-balancer-id] [health-check-id] [flags]
```

## Product catalog:
- Flags:
- --cli.list-links enum[=table]         List all available links for this command (one of "json", "table" or "yaml")
- --description string                  Description
- --force-authentication                Force authentication by sending the header even if this API doesn't require it
- --health-check-id string              health_check_id: ID of the health check you wanna upload (required)
- --healthy-status-code integer         Healthy Status Code (default 200)
- --healthy-threshold-count integer     Healthy Threshold Count (default 8)
- -h, --help                                help for replace
- --initial-delay-seconds integer       Initial Delay Seconds (default 30)
- --interval-seconds integer            Interval Seconds (default 30)
- --load-balancer-id string             load_balancer_id: ID of the attached Load balancer (required)
- --name string                         The Health Check unique name (required)
- --path string                         Path
- --port integer                        Port (required)
- --protocol enum                       HealthCheckProtocol (one of "http" or "tcp") (required)
- --timeout-seconds integer             Timeout Seconds (default 10)
- --unhealthy-threshold-count integer   Unhealthy Threshold Count (default 3)
- -v, --version                             version for replace

## Other commands:
- Global Flags:
- --api-key string           Use your API key to authenticate with the API
- -U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
- use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
- a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
- -t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
- Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
- --debug                    Display detailed log information at the debug level
- --env enum                 Environment to use (one of "pre-prod" or "prod") (default "prod")
- --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
- -o, --output string            Change the output format. Use '--output=help' to know more details.
- -r, --raw                      Output raw data, without any formatting or coloring
- --server-url uri           Manually specify the server to use

## Flags:
```bash

```

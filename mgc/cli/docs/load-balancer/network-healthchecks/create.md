# Create

Create Health Check

## Usage:
```
mgc load-balancer network-healthchecks create [load-balancer-id] [flags]
```

## Flags:
```
    --cli.list-links enum[=table]         List all available links for this command (one of "json", "table" or "yaml")
    --description string                  A brief description of the Health Check
    --healthy-status-code integer         Healthy Status Code: The HTTP status code indicating a healthy response. By default the status is set to 200
    --healthy-threshold-count integer     Healthy Threshold Count: The number of consecutive successful checks before considering the target healthy
-h, --help                                help for create
    --initial-delay-seconds integer       Initial Delay Seconds: The initial delay in seconds before starting Health Checks
    --interval-seconds integer            Interval Seconds: The interval in seconds between Health Checks
    --load-balancer-id uuid               load_balancer_id: ID of the attached Load Balancer (required)
    --name string                         The unique name of the Health Check (max character count: 64) (required)
    --path string                         The path to check for HTTP protocol; ignored for other protocols
    --port integer                        The port number on which the Health Check will be performed (required)
    --protocol enum                       The protocol used for the Health Check (e.g., HTTP, TCP) (one of "http" or "tcp") (required)
    --timeout-seconds integer             Timeout Seconds: The timeout in seconds for each Health Check
    --unhealthy-threshold-count integer   Unhealthy Threshold Count: The number of consecutive failed checks before considering the target unhealthy
-v, --version                             version for create
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


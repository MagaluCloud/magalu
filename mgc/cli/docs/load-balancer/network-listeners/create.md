---
sidebar_position: 2
---
# Create

Create Listener

## Usage:
```
mgc load-balancer network-listeners create [load-balancer-id] [flags]
```

## Flags:
```
    --backend-id uuid               backend_id: ID of the backend that will be the upstream (required)
    --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
    --description string            A brief description of the listener (max character count: 128)
-h, --help                          help for create
    --load-balancer-id uuid         load_balancer_id: ID of the attached Load Balancer (required)
    --name string                   The unique name of the listener (max character count: 64) (required)
    --port integer                  The port number on which the listener will be active (required)
    --protocol enum                 The protocol used by the listener (e.g., TCP, TLS) (one of "tcp" or "tls") (required)
    --tls-certificate-id string     Tls Certificate Id: The ID of the TLS certificate
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


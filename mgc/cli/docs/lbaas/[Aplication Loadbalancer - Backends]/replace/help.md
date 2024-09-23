# Update Backend by ID

## Usage:
```bash
Usage:
  ./mgc lbaas [aplication-loadbalancer-backends] replace [load-balancer-id] [backend-id] [flags]
```

## Product catalog:
- Examples:
- ./mgc lbaas [aplication-loadbalancer-backends] replace --description="Some optional backend description 1" --health-check-name="alb-health-check-1" --name="alb-backend-1"

## Other commands:
- Flags:
- --backend-id integer            Backend Id (required)
- --balance-algorithm enum        BackendBalanceAlgorithm (one of "least_connections", "round_robin" or "source_ip_hash") (required)
- --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
- --description string            Description
- --force-authentication          Force authentication by sending the header even if this API doesn't require it
- --health-check-name string      Health Check Name (required)
- -h, --help                          help for replace
- --load-balancer-id integer      Load Balancer Id (required)
- --name string                   Name (required)
- --parent-id uuid                parent_id: Parent ID to update a backend by ID (required)
- --project-type enum             ProjectType: Project type to update a backend by ID (one of "dbaas", "default" or "k8saas") (required)
- --targets array                 Targets (at least one of: array or array)
- Use --targets=help for more details (required)
- --targets-type enum             BackendType (one of "instance" or "raw") (required)
- -v, --version                       version for replace

## Flags:
```bash
Global Flags:
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
      --server-url uri           Manually specify the server to use
```


# Update Listener by ID

## Usage:
```bash
Usage:
  ./mgc lbaas [network-loadbalancer-listeners] replace [load-balancer-id] [listener-id] [flags]
```

## Product catalog:
- Examples:
- ./mgc lbaas [network-loadbalancer-listeners] replace --backend-name="nlb-backend-1" --description="Some optional listener description 1" --name="nlb-listener-1" --port=80 --public-ip-id="00000000-0000-0000-0000-000000000000" --tls-certificate-name="nlb-tls-certificate-1"

## Other commands:
- Flags:
- --backend-name string           Backend Name (required)
- --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
- --description string            Description
- --force-authentication          Force authentication by sending the header even if this API doesn't require it
- -h, --help                          help for replace
- --listener-id integer           Listener Id (required)
- --load-balancer-id integer      Load Balancer Id (required)
- --name string                   Name (required)
- --parent-id uuid                parent_id: Parent ID to update a listener by ID (required)
- --port integer                  Port (required)
- --project-type enum             ProjectType: Project type to update a listener by ID (one of "dbaas", "default" or "k8saas") (required)
- --protocol enum                 NetworkLoadBalancerProtocol (one of "tcp", "tls" or "udp") (required)
- --public-ip-id string           Public Ip Id
- --tls-certificate-name string   Tls Certificate Name
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


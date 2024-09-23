# Update Load Balancer by ID

## Usage:
```bash
Usage:
  ./mgc lbaas [network-loadbalancer-load-balancers] replace [load-balancer-id] [flags]
```

## Product catalog:
- Examples:
- ./mgc lbaas [network-loadbalancer-load-balancers] replace --description="Some optional load balancer description" --name="nlb-1" --security-groups-ids='["00000000-0000-0000-0000-000000000000"]' --subnet-pool-id="00000000-0000-0000-0000-000000000000" --vpc-id="00000000-0000-0000-0000-000000000000"

## Other commands:
- Flags:
- --acls array(object)                  Acls
- Use --acls=help for more details (required)
- --backends array(object)              Backends
- Use --backends=help for more details (required)
- --cli.list-links enum[=table]         List all available links for this command (one of "json", "table" or "yaml")
- --description string                  Description
- --force-authentication                Force authentication by sending the header even if this API doesn't require it
- --health-checks array(object)         Health Checks
- Use --health-checks=help for more details (required)
- -h, --help                                help for replace
- --listeners array(object)             Listeners
- Use --listeners=help for more details (required)
- --load-balancer-id integer            Load Balancer Id (required)
- --name string                         Name (required)
- --parent-id uuid                      parent_id: Parent ID to update a load balancer by ID (required)
- --project-type enum                   ProjectType: Project type to update a load balancer by ID (one of "dbaas", "default" or "k8saas") (required)
- --security-groups-ids array(string)   Security Groups Ids
- --subnet-pool-id string               Subnet Pool Id
- --tls-certificates array(object)      Tls Certificates
- Use --tls-certificates=help for more details (required)
- --type enum                           NetworkLoadBalancerType (must be "proxy") (required)
- -v, --version                             version for replace
- --visibility enum                     LoadBalancerVisibility (one of "external" or "internal") (required)
- --vpc-id string                       Vpc Id (required)

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


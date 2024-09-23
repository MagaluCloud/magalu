# Create TLS Certificate

## Usage:
```bash
Usage:
  ./mgc lbaas [network-loadbalancer-tls-certificates] create [load-balancer-id] [flags]
```

## Product catalog:
- Examples:
- ./mgc lbaas [network-loadbalancer-tls-certificates] create --certificate="-----BEGIN CERTIFICATE-----" --description="Some optional tls-certificate description 1" --name="nlb-tls-certificate-1" --private-key="-----BEGIN PRIVATE KEY-----"

## Other commands:
- Flags:
- --certificate string            Certificate (required)
- --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
- --description string            Description
- -h, --help                          help for create
- --load-balancer-id integer      Load Balancer Id (required)
- --name string                   The TLS Certificate unique name (required)
- --parent-id uuid                parent_id: Parent ID to create a TLS certificate (required)
- --private-key string            Private Key (required)
- --project-type enum             ProjectType: Project type to create a TLS certificate (one of "dbaas", "default" or "k8saas") (required)
- -v, --version                       version for create

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


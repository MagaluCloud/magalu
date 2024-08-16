# Create a Port with provided vpc_id and x-tenant-id. You can provide a list of security_groups_id or subnets

## Usage:
```bash
Usage:
  ./mgc network port vpcs-ports create [vpc-id] [flags]
```

## Product catalog:
- Flags:
- --cli.list-links enum[=table]        List all available links for this command (one of "json", "table" or "yaml")
- --has-pip                            Has Pip (default true)
- --has-sg                             Has Sg (default true)
- -h, --help                               help for create
- --name string                        Name (between 5 and 100 characters) (required)
- --security-groups-id array(string)   Security Groups Id (default [])
- --subnets array(string)              Subnets (default [])
- -v, --version                            version for create
- --vpc-id string                      vpc_id: ID of the VPC to create port (required)

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
- -o, --output string            Change the ouput format. Use '--output=help' to know more details. (default "yaml")
- -r, --raw                      Output raw data, without any formatting or coloring
- --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri           Manually specify the server to use
- --x-zone string            X-Zone

## Flags:
```bash

```


# VPC Api Product

## Usage:
```bash
Usage:
  mgc network [flags]
  mgc network [command]
```

## Product catalog:
- Commands:
- nat-gateways    Operations related to Nat Gateway
- ports           Operations related to Ports
- public-ips      Operations related to Public IPs
- rules           Operations related to Rules
- security-groups Operations related to Security Groups
- subnetpools     Operations related to Subnet Pools
- subnets         Operations related to Subnets
- vpcs            Operations related to VPCs

## Other commands:
- Flags:
- -h, --help      help for network
- -v, --version   version for network

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
      --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
  -o, --output string            Change the output format. Use '--output=help' to know more details.
  -r, --raw                      Output raw data, without any formatting or coloring
```


# Lbaas RP API

## Usage:
```bash
Usage:
  ./mgc lbaas [flags]
  ./mgc lbaas [command]
```

## Product catalog:
- Commands:
- application-backends      Application Load Balancer Backends (Target Pools)
- application-certificates  Application Load Balancer Tls Certificates
- application-healthchecks  Application Load Balancer HealthCheks
- application-listeners     Application Load Balancer Listeners
- application-loadbalancers Application Load Balancer
- network-backends          Network Load Balancer Backends (Target Pools)
- network-certificates      Network Load Balancer Tls Certificates
- network-healthchecks      Network Load Balancer HealthCheks
- network-listeners         Network Load Balancer Listeners
- network-loadbalancers     Network Load Balancer

## Other commands:
- Flags:
- -h, --help      help for lbaas
- -v, --version   version for lbaas

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

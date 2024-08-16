# Start conciliator validation

## Usage:
```bash
Usage:
  ./mgc network backoffice-conciliator create [flags]
```

## Product catalog:
- Flags:
- -h, --help      help for create
- -v, --version   version for create

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
- --limit integer            Limit of tenants
- --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
- -o, --output string            Change the ouput format. Use '--output=help' to know more details. (default "yaml")
- -r, --raw                      Output raw data, without any formatting or coloring
- --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --resource-type enum       Resource Type: Resources to start the conciliator (one of "public_ip", "rule", "security_group" or "security_group, rule, public_ip") (default "security_group, rule, public_ip")
- --server-url uri           Manually specify the server to use
- --skip integer             Number of tenants to skip

## Flags:
```bash

```


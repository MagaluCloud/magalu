# List all Parameter Groups for a x-tenant-id

## Usage:
```bash
Usage:
  mgc dbaas parameter-groups list [flags]
```

## Product catalog:
- Examples:
- mgc dbaas parameter-groups list --type="USER"

## Other commands:
- Flags:
- --control.limit integer    The maximum number of items per page. (range: 1 - 25) (default 10)
- --control.offset integer   The number of items to skip before starting to collect the result set. (min: 0)
- --engine-id uuid           Value referring to engine Id.
- -h, --help                     help for list
- --type enum                Type (one of "SYSTEM" or "USER")
- -v, --version                  version for list

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
      --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri           Manually specify the server to use
```


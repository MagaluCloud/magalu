# Policy-related commands

## Usage:
```bash
Usage:
  ./mgc object-storage buckets policy [flags]
  ./mgc object-storage buckets policy [command]
```

## Product catalog:
- Commands:
- delete      Delete policy document for the specified bucket
- get         Get the policy document for the specified bucket
- set         Set policy document for the specified bucket

## Other commands:
- Flags:
- -h, --help   help for policy

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
      --host string              URL to override the default host. Ex. https://api.magalu.com.br or http://localhost/v1/route
      --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
  -o, --output string            Change the output format. Use '--output=help' to know more details.
  -r, --raw                      Output raw data, without any formatting or coloring
```


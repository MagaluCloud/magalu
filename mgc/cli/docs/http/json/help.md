# JSON HTTP access

## Usage:
```bash
Usage:
  mgc http json [flags]
  mgc http json [command]
```

## Product catalog:
- Commands:
- delete      Call using HTTP DELETE
- get         Call using HTTP GET
- head        Call using HTTP HEAD
- options     Call using HTTP OPTIONS
- patch       Call using HTTP PATCH
- post        Call using HTTP POST
- put         Call using HTTP PUT
- trace       Call using HTTP TRACE

## Other commands:
- Flags:
- -h, --help   help for json

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


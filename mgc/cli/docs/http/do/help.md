# Execute generic http request

## Usage:
```bash
Usage:
  ./mgc http do [flags]
```

## Product catalog:
- Flags:
- --body string
- --headers object           Use --headers=help for more details
- -h, --help                     help for do
- --method enum              (one of "DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT" or "TRACE") (required)
- --security array(object)   Use --security=help for more details
- --url string               Golang template with the URL (required)

## Other commands:
- Global Flags:
- --api-key string           Use your API key to authenticate with the API
- -U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
- use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
- a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
- -t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
- Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
- --debug                    Display detailed log information at the debug level
- --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
- -o, --output string            Change the output format. Use '--output=help' to know more details. (default "yaml")
- -r, --raw                      Output raw data, without any formatting or coloring
- --region enum              Region to reach the service (one of "br-mgl-1", "br-ne-1" or "br-se-1") (default "br-ne-1")
- --server-url uri           Manually specify the server to use

## Flags:
```bash

```


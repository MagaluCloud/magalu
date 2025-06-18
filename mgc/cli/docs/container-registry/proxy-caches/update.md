# Update

Update a proxycache by uuid.

## Usage:
```
mgc container-registry proxy-caches update [proxy-cache-id] [flags]
```

## Flags:
```
    --access-key string             A string consistent with provider access_id.
    --access-secret string          A string consistent with provider access_secret.
    --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
    --description string            A string.
-h, --help                          help for update
    --name string                   A unique name for each tenant, used for the proxy-cache. It must be written in lowercase letters and consists only of numbers and letters, up to a limit of 63 characters.
    --proxy-cache-id uuid           Proxy cache's UUID. (required)
    --url string                    An Endpoint URL for the proxied registry. Example URL for available providers can be checked through mcr-api or mgccli.
```

## Global Flags:
```
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


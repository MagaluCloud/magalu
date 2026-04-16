---
sidebar_position: 5
---
# Vulnerabilities

Returns paginated vulnerabilities found in a specific image scan.

## Usage:
```
mgc container-registry scans vulnerabilities [scan-id] [flags]
```

## Examples:
```
mgc container-registry scans vulnerabilities --severity='["high"]'
```

## Flags:
```
    --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
    --control.limit integer         Limit. (min: 1)
    --control.offset integer        Offset. (min: 0)
    --control.sort string           Fields to use as reference to sort. (pattern: (^[\w-]+:(asc|desc)(,[\w-]+:(asc|desc))*)?$)
    --cve-id string                 Filter vulnerabilities by CVE ID.
    --fixable                       Filter vulnerabilities that have a fixed version available.
-h, --help                          help for vulnerabilities
    --package-name string           Filter vulnerabilities by package name.
    --scan-id uuid                  Scan UUID. (required)
    --severity array(enum)          Filter vulnerabilities by severity.
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
-o, --output string            Change the output format. You can use 'yaml', 'json' or 'table'.
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
```


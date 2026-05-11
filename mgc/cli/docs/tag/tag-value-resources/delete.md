---
sidebar_position: 3
---
# Delete

Delete a link between a tag-value and a cloud resource

## Usage:
```
mgc tag tag-value-resources delete [tag-name] [value-name] [resource-id] [flags]
```

## Examples:
```
mgc tag tag-value-resources delete --tag-name="kubernetes-expenses" --value-name="test-labs"
```

## Flags:
```
    --force-authentication   Force authentication by sending the header even if this API doesn't require it
-h, --help                   help for delete
    --resource-id uuid       Resource Id (required)
    --tag-name string        Tag Name: name for tag (required)
    --value-name string      tag value name, is allowed only one name per tag_id (required)
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
    --server-url uri           Manually specify the server to use
```


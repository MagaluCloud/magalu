---
sidebar_position: 1
---
# List

List the resources you have access to, with the tags attached to each one.

## Usage:
```
mgc tag resources list [flags]
```

## Examples:
```
mgc tag resources list --external-id="31201f93-f5f4-4cf1-ba9c-bfed0717f4ac"
```

## Flags:
```
    --control.limit integer     Maximum number of items to return per page (range: 1 - 100)
    --control.offset integer    Number of items to skip before the first result
    --control.sort string       Sort criteria for the result
    --external-id string        External ID of the resource to filter by
-h, --help                      help for list
    --region enum               Region to filter by (one of "br-ne1", "br-se1" or "global")
    --resource-type-name enum   Resource type to filter by, prefixed by its product (one of "bs.snapshot", "bs.volume", "cr.registry", "cr.repository", "db.cluster", "db.instance", "db.parameter-group", "db.replica", "db.snapshot", "k8s.cluster", "k8s.nodepool", "lb.network-acl", "lb.network-backend", "lb.network-certificate", "lb.network-healthcheck", "lb.network-listener", "lb.network-loadbalancer", "net.nat-gateway", "net.port", "net.public-ip", "net.rule", "net.security-group", "net.subnet", "net.vpc", "os.bucket", "os.object", "un.unknown", "vm.image", "vm.instance" or "vm.snapshot")
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


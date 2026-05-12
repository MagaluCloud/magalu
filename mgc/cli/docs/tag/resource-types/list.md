---
sidebar_position: 1
---
# List

List all resource types

## Usage:
```
mgc tag resource-types list [flags]
```

## Flags:
```
    --control.limit integer    Limit: Number of items per page (range: 1 - 100)
    --control.offset integer   Offset for pagination (min: 0)
    --control.sort string      Sorting criteria
-h, --help                     help for list
    --name enum                ResourceEnum: resource type name, must be related to a product (one of "bs.snapshot", "bs.volume", "cr.registry", "cr.repository", "db.cluster", "db.instance", "db.parameter-group", "db.replica", "db.snapshot", "k8s.cluster", "k8s.nodepool", "lb.network-acl", "lb.network-backend", "lb.network-certificate", "lb.network-healthcheck", "lb.network-listener", "lb.network-loadbalancer", "net.nat-gateway", "net.port", "net.public-ip", "net.rule", "net.security-group", "net.subnet", "net.vpc", "os.bucket", "os.object", "vm.image", "vm.instance" or "vm.snapshot")
    --product enum             ProductEnum: product owner of a resource (one of "block-storage", "container-registry", "database", "kubernetes", "load-balancer", "network", "object-storage" or "virtual-machine")
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


---
sidebar_position: 0
---
# Commands-Reference

Magalu Cloud CLI is a command-line interface for the Magalu Cloud.
It allows you to interact with the Magalu Cloud to manage your resources.

## Usage:
```
mgc [flags]
mgc [command]
```

## Product catalog:
```
audit              Cloud Events API Product.
block-storage      Block Storage API Product
container-registry Magalu Container Registry product API.
dbaas              DBaaS API Product.
kubernetes         APIs related to the Kubernetes product.
load-balancer      Lbaas API: create and manage Load Balancers
network            VPC Api Product
object-storage     Operations for Object Storage
virtual-machine    Virtual Machine Api Product
```

## Other commands:
```
completion         Generate the autocompletion script for the specified shell
help               Help about any command
```

## Flags:
```
    --api-key string           Use your API key to authenticate with the API
-U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                               use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                               a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
-t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                               Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
    --debug                    Display detailed log information at the debug level
-h, --help                     help for mgc
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. Use '--output=help' to know more details.
-r, --raw                      Output raw data, without any formatting or coloring
-v, --version                  version for mgc
```

## Settings:
```
auth               Actions with ID Magalu to log in, API Keys, refresh tokens, change tenants and others
config             Manage CLI Configuration values
profile            Manage account settings, including SSH keys and related configurations
workspace          Manage workspaces for isolated auth and config settings
```


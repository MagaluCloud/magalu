---
sidebar_position: 0
---
# Federations

Domains Federations of organization. | federations

## Usage:
```
mgc iam domains federations [flags]
mgc iam domains federations [command]
```

## Commands:
```
list        List Federation
oidc        Create Domain Federation Oidc
saml        Create Domain Federation Saml
update      Update Domain Federation Resource
```

## Flags:
```
-h, --help   help for federations
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
    --iam-project-id string    Project scope for IAM commands only. Overrides the configured IAM project
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. You can use 'yaml', 'json' or 'table'.
    --project-id string        Project to scope the requests. Overrides the configured project
-r, --raw                      Output raw data, without any formatting or coloring
```


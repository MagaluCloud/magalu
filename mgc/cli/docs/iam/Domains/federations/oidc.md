---
sidebar_position: 3
---
# Oidc

Create Domain Federation Oidc

## Usage:
```
mgc iam domains federations oidc [domain-id] [flags]
```

## Flags:
```
    --authorization-endpoint string   Authorization Endpoint (required)
    --cli.list-links enum[=table]     List all available links for this command (one of "json", "table" or "yaml")
    --client-id string                Client Id (required)
    --client-secret string            Client Secret (required)
    --domain-id uuid                  Domain Id (required)
-h, --help                            help for oidc
    --jwks-uri string                 Jwks Uri (required)
    --name string                     Name (required)
    --project-id string               Project to scope this command to. Overrides the configured project for this invocation
    --scope enum                      Scope this command applies to: 'default' for the tenant's default project, 'tenant' for the entire tenant. The IAM API encodes the default project as the tenant id (one of "default" or "tenant")
    --scopes array(string)            Scopes (required)
    --token-endpoint string           Token Endpoint (required)
    --userinfo-endpoint string        Userinfo Endpoint (required)
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
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. You can use 'yaml', 'json' or 'table'.
-r, --raw                      Output raw data, without any formatting or coloring
    --server-url uri           Manually specify the server to use
```


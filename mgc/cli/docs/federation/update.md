---
sidebar_position: 5
---
# Update

Edit Federation Resource

## Usage:
```
mgc federation update [domain-id] [federation-id] [flags]
```

## Flags:
```
    --domain-id uuid                       Domain Id (required)
    --federation-id uuid                   Federation Id (required)
-h, --help                                 help for update
    --oidc object                          (properties: authorization_endpoint, client_id, client_secret, jwks_uri, name, scopes, token_endpoint and userinfo_endpoint)
                                           Use --oidc=help for more details
    --oidc.authorization-endpoint string   oidc's authorization_endpoint property: Authorization Endpoint
                                           This is the same as '--oidc=authorization_endpoint:string'.
    --oidc.client-id string                oidc's client_id property: Client Id
                                           This is the same as '--oidc=client_id:string'.
    --oidc.client-secret string            oidc's client_secret property: Client Secret
                                           This is the same as '--oidc=client_secret:string'.
    --oidc.jwks-uri string                 oidc's jwks_uri property: Jwks Uri
                                           This is the same as '--oidc=jwks_uri:string'.
    --oidc.name string                     oidc's name property: Name
                                           This is the same as '--oidc=name:string'.
    --oidc.scopes array                    oidc's scopes property: Scopes
                                           This is the same as '--oidc=scopes:array'.
    --oidc.token-endpoint string           oidc's token_endpoint property: Token Endpoint
                                           This is the same as '--oidc=token_endpoint:string'.
    --oidc.userinfo-endpoint string        oidc's userinfo_endpoint property: Userinfo Endpoint
                                           This is the same as '--oidc=userinfo_endpoint:string'.
    --saml object                          (properties: entity_id, idp_certificate, name and sso_url)
                                           Use --saml=help for more details
    --saml.entity-id string                saml's entity_id property: Entity Id
                                           This is the same as '--saml=entity_id:string'.
    --saml.idp-certificate string          saml's idp_certificate property: Idp Certificate
                                           This is the same as '--saml=idp_certificate:string'.
    --saml.name string                     saml's name property: Name
                                           This is the same as '--saml=name:string'.
    --saml.sso-url string                  saml's sso_url property: Sso Url
                                           This is the same as '--saml=sso_url:string'.
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


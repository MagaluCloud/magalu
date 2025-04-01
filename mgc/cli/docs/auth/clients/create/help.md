# Create

Create new client (Oauth Application)

## Usage:
```
mgc auth clients create [name] [description] [redirect-uris] [backchannel-logout-session] [client-term-url] [client-privacy-term-url] [audiences] [email] [request-reason] [icon] [access-token-expiration] [always-require-login] [backchannel-logout-uri] [oidc-audience] [refresh-token-custom-expires-enabled] [refresh-token-exp] [support-url] [grant-types] [flags]
```

## Examples:
```
mgc auth clients create --access-token-expiration=7200 --audiences="public" --description="Client description" --email="client@example.com" --name="Client Name" --oidc-audience="public" --refresh-token-exp=15778476
```

## Flags:
```
    --access-token-expiration integer        Access token expiration (in seconds)
    --always-require-login                   Must ignore active Magalu ID session and always require login
    --audiences string                       Client audiences (separated by space)
    --backchannel-logout-session             Client requires backchannel logout session
    --backchannel-logout-uri string          Backchannel logout URI
    --client-privacy-term-url string         URL to privacy term (required)
    --client-term-url string                 URL to terms of use (required)
    --description string                     Description of new client (required)
    --email string                           Email of new client
    --grant-types string                     Grant types the client can request for token generation (separated by space)
-h, --help                                   help for create
    --icon string                            URL for client icon
    --name string                            Name of new client (required)
    --oidc-audience string                   OIDC audience (separated by space)
    --redirect-uris string                   Redirect URIs (separated by space) (required)
    --refresh-token-custom-expires-enabled   Use custom value for refresh token expiration
    --refresh-token-exp integer              Custom refresh token expiration value (in seconds)
    --request-reason string                  Note to inform the reason for creating the client. Will help with the application approval process
    --support-url string                     Support URL
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
-o, --output string            Change the output format. Use '--output=help' to know more details.
-r, --raw                      Output raw data, without any formatting or coloring
```


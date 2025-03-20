# Update a client (Oauth Application)

## Usage:
```bash
Usage:
  mgc auth clients update [id] [name] [description] [redirect-uris] [icon] [access-token-expiration] [always-require-login] [client-privacy-term-url] [client-term-url] [audiences] [backchannel-logout-session] [backchannel-logout-uri] [oidc-audience] [refresh-token-custom-expires-enabled] [refresh-token-exp] [request-reason] [support-url] [flags]
```

## Product catalog:
- Examples:
- mgc auth clients update --access-token-expiration=7200 --audiences="public" --description="Client description" --name="Client Name" --refresh-token-exp=15778476

## Other commands:
- Flags:
- --access-token-expiration integer        Access token expiration (in seconds)
- --always-require-login                   Must ignore active Magalu ID session and always require login
- --audiences string                       Client audiences (separated by space)
- --backchannel-logout-session             Client requires backchannel logout session
- --backchannel-logout-uri string          Backchannel logout URI
- --client-privacy-term-url string         URL to privacy term
- --client-term-url string                 URL to terms of use
- --description string                     Description of new client
- -h, --help                                   help for update
- --icon string                            URL for client icon
- --id string                              UUID of client (required)
- --name string                            Name of new client
- --oidc-audience string                   Audiences for ID token
- --redirect-uris string                   Redirect URIs (separated by space)
- --refresh-token-custom-expires-enabled   Use custom value for refresh token expiration
- --refresh-token-exp integer              Custom refresh token expiration value (in seconds)
- --request-reason string                  Note to inform the reason for creating the client. Will help with the application approval process
- --support-url string                     URL for client support

## Flags:
```bash
Global Flags:
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


# Create new client (Oauth Application)

## Usage:
```bash
Usage:
  ./cli auth clients create [name] [description] [redirect-uris] [icon] [access-token-expiration] [always-require-login] [client-privacy-term-url] [terms-of-use] [audiences] [backchannel-logout-session-enabled] [backchannel-logout-uri] [oidc-audiences] [refresh-token-custom-expires-enabled] [refresh-token-expiration] [reason] [flags]
```

## Product catalog:
- Examples:
- ./cli auth clients create --access-token-expiration=7200 --audiences="public" --description="Client description" --name="Client Name" --refresh-token-expiration=15778476

## Other commands:
- Flags:
- --access-token-expiration integer        Access token expiration (in seconds)
- --always-require-login                   Must ignore active Magalu ID session and always require login
- --audiences string                       Client audiences (separated by space)
- --backchannel-logout-session-enabled     Client requires backchannel logout session
- --backchannel-logout-uri string          Backchannel logout URI
- --client-privacy-term-url string         URL to privacy term (required)
- --description string                     Description of new client (required)
- -h, --help                                   help for create
- --icon string                            URL for client icon
- --name string                            Name of new client (required)
- --oidc-audiences string                  Audiences for ID token
- --reason string                          Note to inform the reason for creating the client. Will help with the application approval process (required)
- --redirect-uris string                   Redirect URIs (separated by space) (required)
- --refresh-token-custom-expires-enabled   Use custom value for refresh token expiration
- --refresh-token-expiration integer       Custom refresh token expiration value (in seconds)
- --terms-of-use string                    URL to terms of use (required)

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
```

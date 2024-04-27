# Retrieve a list of Volumes for the currently authenticated tenant.

## Usage:
```bash
#### Notes
- Use the expand argument to obtain additional details about the Volume Type.
```

## Product catalog:
- Usage:
- ./cli block-storage volumes list [flags]

## Other commands:
- Flags:
- --control.limit integer     Limit (min: 1) (default 50)
- --control.offset integer    Offset (min: 0)
- --control.sort string       Sort (pattern: ^(^[\w-]+:(asc|desc)(,[\w-]+:(asc|desc))*)?$) (default "created_at:asc")
- --expand array(string)     Expand: You can get more detailed info about: ['volume_type', 'attachment'] (default [])
- -h, --help                     help for list
- -v, --version                  version for list

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-ne1")
      --server-url uri         Manually specify the server to use
```

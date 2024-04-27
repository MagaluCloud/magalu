# Tenants work like sub-accounts. You may have more than one Tenant under your
Magalu Cloud account and they each store their data separately, but are billed
under the same account

## Usage:
```bash
Usage:
  ./cli auth tenant [flags]
  ./cli auth tenant [command]
```

## Product catalog:
- Commands:
- current     Get the currently active Tenant
- list        List all available tenants for current login
- set         Set the active Tenant to be used for all subsequential requests

## Other commands:
- Additional Commands:
- select      call "list", prompt selection and then "set"

## Flags:
```bash
Flags:
  -h, --help   help for tenant
```


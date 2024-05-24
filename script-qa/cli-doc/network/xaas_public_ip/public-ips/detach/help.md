# Detach a Public IP to a Port

## Usage:
```bash
Usage:
  ./mgc network xaas-public-ip public-ips detach [public-ip-id] [port-id] [flags]
```

## Product catalog:
- Flags:
- --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
- -h, --help                          help for detach
- --port-id string                Port ID: Id of the Port to detach the Public IP (required)
- --project-type enum             project_type: Project type to delete tenant's public ip (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- --public-ip-id string           Public IP ID: Id of the Public IP to detach port to (required)
- -v, --version                       version for detach

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri         Manually specify the server to use

## Flags:
```bash

```


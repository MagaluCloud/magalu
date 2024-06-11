# XAAS detach a Security Group to a Port with provided port_id, security_group_id, x-tenant-id of an specific project type

## Usage:
```bash
Usage:
  ./mgc network xaas-port ports detach [port-id] [security-group-id] [flags]
```

## Product catalog:
- Flags:
- --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
- -h, --help                          help for detach
- --port-id string                port_id: ID of the Port to detach security group (required)
- --project-type enum             project_type: Project type to detach security group (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- --security-group-id string      security_group_id: ID of the Security Group to detach (required)
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

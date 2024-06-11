# Create a Security Group async, returning its ID. To monitor the creation progress, please check the status in the service message or implement polling.

## Usage:
```bash
Usage:
  ./mgc network xaas-security-group security-groups create [flags]
```

## Product catalog:
- Flags:
- --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
- --description string            Description
- -h, --help                          help for create
- --name string                   Name (between 5 and 100 characters) (required)
- --project-type enum             project_type: Project type to create Security Group (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- -v, --version                       version for create

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri         Manually specify the server to use

## Flags:
```bash

```

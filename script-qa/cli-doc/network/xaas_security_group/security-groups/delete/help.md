# Delete a security group from the provided tenant_id

## Usage:
```bash
Usage:
  ./mgc network xaas-security-group security-groups delete [security-group-id] [flags]
```

## Product catalog:
- Flags:
- -h, --help                       help for delete
- --project-type enum          project_type: Project type to delete Security Group (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- --security-group-id string   securityGroupId: Id of the Security Group to delete (required)
- -v, --version                    version for delete

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri         Manually specify the server to use

## Flags:
```bash

```

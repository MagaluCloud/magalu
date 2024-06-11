# Returns a list of XAAS subnets from the provided vpc_id

## Usage:
```bash
Usage:
  ./mgc network xaas-subnets list [vpc-id] [flags]
```

## Product catalog:
- Flags:
- -h, --help                help for list
- --project-type enum   project_type: Project type to list Subnets (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- -v, --version             version for list
- --vpc-id string       VPC Id: Id of the VPC to list Subnets (required)

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri         Manually specify the server to use

## Flags:
```bash

```

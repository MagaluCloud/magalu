# Create a Port on an xaas project with provided vpc_id and x-tenant-id. You can provide a list of security_groups_id or subnets

## Usage:
```bash
Usage:
  ./mgc network xaas-port create [vpc-id] [flags]
```

## Product catalog:
- Flags:
- --has-pip                            Has Pip
- --has-sg                             Has Sg
- -h, --help                               help for create
- --name string                        Name (between 5 and 100 characters) (required)
- --project-type enum                  project_type: Project type to create port (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- --security-groups-id array(string)   Security Groups Id
- --subnets array(string)              Subnets
- -v, --version                            version for create
- --vpc-id string                      vpc_id: ID of the VPC to create port (required)

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri         Manually specify the server to use
- --x-zone string          X-Zone

## Flags:
```bash

```


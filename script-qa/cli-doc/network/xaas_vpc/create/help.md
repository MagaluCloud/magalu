# Create a VPC for a provided tenant_id of an project type

## Usage:
```bash
Usage:
  ./mgc network xaas-vpc create [flags]
```

## Product catalog:
- Flags:
- --cidr string            CIDR from which VPC's subnets will get IPs            Can be used instead of 'subnetpool_id' field.
- --description string     Description
- -h, --help                   help for create
- --name string            Name (required)
- --project-type enum      project_type: Project type to create VPC (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- --subnetpool-id string   Subnetpool Id: ID of SubnetPool             from which the VPC will get CIDRs.             Can be used instead of 'cidr' field.
- -v, --version                version for create

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri         Manually specify the server to use

## Flags:
```bash

```


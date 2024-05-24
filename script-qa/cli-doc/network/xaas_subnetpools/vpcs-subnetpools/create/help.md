# Create a Subnet Pool

## Usage:
```bash
Usage:
  ./mgc network xaas-subnetpools vpcs-subnetpools create [vpc-id] [flags]
```

## Product catalog:
- Flags:
- --address-scope-id string         Address Scope Id: The ID of the address scope for the subnet pool
- --cli.list-links enum[=table]     List all available links for this command (one of "json", "table" or "yaml")
- --default-prefix-length integer   The default prefix length for a subnet in the pool. (default 26)
- --description string              The description for the subnet pool (required)
- -h, --help                            help for create
- --max-prefix-length integer       Max Prefix Length: The maximum prefix length for a subnet in the pool. (default 28)
- --min-prefix-length integer       Min Prefix Length: The minimum prefix length for a subnet in the pool. (default 24)
- --name string                     The name of the subnet pool. (required)
- --pool-prefix string              Pool Prefix: The CIDR notation prefix for the subnet pool. (default "172.26.0.0/16")
- --project-type enum               project_type: Project type to create Subnet Pool (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- -v, --version                         version for create
- --vpc-id string                   VPC Id: Id of the VPC to create the Subnet Pool (required)

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri         Manually specify the server to use

## Flags:
```bash

```


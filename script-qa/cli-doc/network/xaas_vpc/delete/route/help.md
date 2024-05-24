# Delete a static route for VPCs primary Router

## Usage:
```bash
Usage:
  ./mgc network xaas-vpc delete route [vpc-id] [flags]
```

## Product catalog:
- Flags:
- --destination string   Destination (required)
- -h, --help                 help for route
- --nexthop string       Nexthop (required)
- --project-type enum    project_type: Project type to delete route (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- -v, --version              version for route
- --vpc-id string        VPC Id: Id of the VPC to delete route (required)

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri         Manually specify the server to use

## Flags:
```bash

```


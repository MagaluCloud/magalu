# Returns a list of ports for a provided vpc_id and x-tenant-id. The list will be paginated, it means you can easily find what you need just setting the page number(_offset) and the quantity of items per page(_limit). The level of detail can also be set

## Usage:
```bash
Usage:
  ./mgc network port vpcs-ports list [vpc-id] [flags]
```

## Product catalog:
- Flags:
- --control.limit integer          Items Per Page (min: 1) (default 10)
- --control.offset integer         Page Number (min: 1) (default 1)
- --detailed                       Detailed (default true)
- -h, --help                           help for list
- --port-id-list array(anyValue)   Port Id List
- Use --port-id-list=help for more details (default [])
- -v, --version                        version for list
- --vpc-id anyValue                vpc_id: ID of VPC to list ports
- Use --vpc-id=help for more details (required)

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text
- --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
- --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
- --server-url uri         Manually specify the server to use

## Flags:
```bash

```


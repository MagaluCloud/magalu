---
sidebar_position: 2
---
# Create

Create a new Virtual Private Cloud (VPC) Peering.

## Usage:
```
mgc network vpcs-peerings create [flags]
```

## Examples:
```
mgc network vpcs-peerings create --description="Conexão entre a VPC de produção e a VPC de banco de dados" --name="peering-prod-to-db" --vpcs.accepter-vpc-id="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" --vpcs.requester-vpc-id="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
```

## Flags:
```
    --description string             Description: Descrição opcional do Peering.
-h, --help                           help for create
    --name string                    Name: Nome do Peering. Deve ser alfanumérico ou hífens, max 50 chars. (required)
    --vpcs object                    VpcPeeringCreateVpcsPayload (properties: accepter_vpc_id and requester_vpc_id)
                                     Use --vpcs=help for more details (required)
    --vpcs.accepter-vpc-id string    VpcPeeringCreateVpcsPayload: Accepter Vpc Id: ID da VPC que recebera o convite de peering (Accepter).
                                     This is the same as '--vpcs=accepter_vpc_id:string'.
    --vpcs.requester-vpc-id string   VpcPeeringCreateVpcsPayload: Requester Vpc Id: ID da VPC que esta solicitando o peering (Requester).
                                     This is the same as '--vpcs=requester_vpc_id:string'.
```

## Global Flags:
```
    --api-key string           Use your API key to authenticate with the API
-U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                               use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                               a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
-t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                               Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
    --debug                    Display detailed log information at the debug level
    --env enum                 Environment to use (one of "pre-prod" or "prod") (default "prod")
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. You can use 'yaml', 'json' or 'table'.
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
```


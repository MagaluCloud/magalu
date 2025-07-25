---
sidebar_position: 7
---
# Restore

Restore a snapshot of an instance with the current tenant which is logged in. 

## Usage:
```
mgc virtual-machine snapshots restore [id] [flags]
```

## Examples:
```
mgc virtual-machine snapshots restore --machine-type.id="9ec75090-2872-4f51-8111-53d05d96d2c6" --machine-type.name="some_resource_name" --network.associate-public-ip=true --network.interface.id="9ec75090-2872-4f51-8111-53d05d96d2c6" --network.interface.security-groups='[{"id":"9ec75090-2872-4f51-8111-53d05d96d2c6"}]' --network.vpc.id="9ec75090-2872-4f51-8111-53d05d96d2c6" --network.vpc.name="some_resource_name"
```

## Flags:
```
    --availability-zone string                  Availability Zone (between 1 and 255 characters)
    --cli.list-links enum[=table]               List all available links for this command (one of "json", "table" or "yaml")
-h, --help                                      help for restore
    --id string                                 Id (required)
    --machine-type object                       Machine Type (at least one of: single property: id or single property: name)
                                                Use --machine-type=help for more details (required)
    --machine-type.id string                    Machine Type: Id (between 1 and 255 characters)
                                                This is the same as '--machine-type=id:string'.
    --machine-type.name string                  Machine Type: Name (between 1 and 255 characters)
                                                This is the same as '--machine-type=name:string'.
    --name string                               Name (between 1 and 255 characters) (required)
    --network object                            (properties: associate_public_ip, interface and vpc)
                                                Use --network=help for more details
    --network.associate-public-ip boolean       network's associate_public_ip property: Associate Public Ip
                                                This is the same as '--network=associate_public_ip:boolean'.
    --network.interface object                  network's interface property: Interface (at least one of: single property: id or single property: security_groups)
                                                Use --network.interface=help for more details
                                                This is the same as '--network=interface:object'.
    --network.interface.id string               Interface: Id (between 1 and 255 characters)
                                                This is the same as '--network.interface=id:string'.
    --network.interface.security-groups array   Interface: Security Groups
                                                Use --network.interface.security-groups=help for more details
                                                This is the same as '--network.interface=security_groups:array'.
    --network.vpc object                        network's vpc property: Vpc (at least one of: single property: id or single property: name)
                                                Use --network.vpc=help for more details
                                                This is the same as '--network=vpc:object'.
    --network.vpc.id string                     Vpc: Id (between 1 and 255 characters)
                                                This is the same as '--network.vpc=id:string'.
    --network.vpc.name string                   Vpc: Name (between 1 and 255 characters)
                                                This is the same as '--network.vpc=name:string'.
    --ssh-key-name string                       Ssh Key Name
    --user-data string                          User Data (between 1 and 65000 characters)
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
-o, --output string            Change the output format. Use '--output=help' to know more details.
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
```


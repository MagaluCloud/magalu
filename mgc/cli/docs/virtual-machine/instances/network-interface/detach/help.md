# Detach a non primary network interface from an instance for a default project

## Usage:
```bash
Usage:
  mgc virtual-machine instances network-interface detach [flags]
```

## Product catalog:
- Examples:
- mgc virtual-machine instances network-interface detach --instance.id="9ec75090-2872-4f51-8111-53d05d96d2c6" --instance.name="some_resource_name" --network.interface.id="9ec75090-2872-4f51-8111-53d05d96d2c6"

## Other commands:
- Flags:
- -h, --help                          help for detach
- --instance object               Instance (at least one of: single property: id or single property: name)
- Use --instance=help for more details (required)
- --instance.id string            Instance: Id (between 1 and 255 characters)
- This is the same as '--instance=id:string'.
- --instance.name string          Instance: Name (between 1 and 255 characters)
- This is the same as '--instance=name:string'.
- --network object                NICRequestInterfaceField (single property: interface)
- Use --network=help for more details (required)
- --network.interface object      NICRequestInterfaceField: ID (single property: id)
- Use --network.interface=help for more details
- This is the same as '--network=interface:object'.
- --network.interface.id string   ID: Id (between 1 and 255 characters)
- This is the same as '--network.interface=id:string'.
- -v, --version                       version for detach

## Flags:
```bash
Global Flags:
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


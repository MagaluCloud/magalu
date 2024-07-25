# Create an Event using the old format, only used for some internal systems, for integration please use: /events/

## Usage:
```bash
Usage:
  ./mgc events internal create [flags]
```

## Product catalog:
- Examples:
- ./mgc events internal create --action-tenant-id="00000000-00000000-00000000-00000000" --action="create" --id="00000000-00000000-00000000-00000000" --product="compute" --project-tenant-id="00000000-00000000-00000000-00000000" --project-type="default" --source="https://api.com/v1/instances/00000000-00000000-00000000-00000000" --time="2024-07-16T22:50:00" --type="instance"

## Other commands:
- Flags:
- --action string            The value describing the action of event related to the event. (required)
- --action-tenant-id uuid    Action Tenant Id: ID of the tenant which requested the change (required)
- --data object              The raw event about the occurrence
- Use --data=help for more details (required)
- -h, --help                     help for create
- --id uuid                  Identifies the event. (required)
- --product string           The identification in which producer product an event occur (required)
- --project-tenant-id uuid   Project Tenant Id: A unique identifier of the principal that triggered the occurrence. (required)
- --project-type string      Project Type: The identification in which producer type an event occur (required)
- --source uri               Source: Identifies the context in which the event occurred. (min character count: 1) (required)
- --specversion string       Specversion: Version of the CloudEvents specification which the event uses. (default "1.0")
- --time date-time           Timestamp of when the occurrence happened. (required)
- --type string              The value describing the type of event related to the originating occurrence. (required)
- -v, --version                  version for create

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri         Manually specify the server to use
```


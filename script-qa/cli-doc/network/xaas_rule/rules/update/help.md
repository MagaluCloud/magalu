# Update async creation of Rule

## Usage:
```bash
Usage:
  ./mgc network xaas-rule rules update [rule-id] [flags]
```

## Product catalog:
- Examples:
- ./mgc network xaas-rule rules update --rules-zones='[{"resource_id":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","zone":"zone_name"}]'

## Other commands:
- Flags:
- --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
- --direction string              Direction
- --error string                  Error
- --external-id string            External Id
- -h, --help                          help for update
- --port-range-max integer        Port Range Max
- --port-range-min integer        Port Range Min
- --project-type enum             project_type: Project type to create rule (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- --protocol string               Protocol
- --remote-group-id string        Remote Group Id
- --remote-ip-prefix string       Remote Ip Prefix
- --resource-id string            Resource Id
- --rule-id string                Rule ID: Id of the Rule (required)
- --rules-zones array(object)     Rules Zones
- Use --rules-zones=help for more details (default [])
- --security-group-id string      Security Group Id
- --status enum                   RuleStatus (one of "created", "error" or "pending") (required)
- -v, --version                       version for update

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri         Manually specify the server to use
      --x-request-id string    X-Request-Id: Request id of Rule to update
```


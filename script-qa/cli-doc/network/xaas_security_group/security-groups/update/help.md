# Update async creation of Security Group

## Usage:
```bash
Usage:
  ./mgc network xaas-security-group security-groups update [security-group-id] [flags]
```

## Product catalog:
- Examples:
- ./mgc network xaas-security-group security-groups update --rules='[{"api_id":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","direction":"egress","error":"null","external_id":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","port_range_max":8028,"port_range_min":8028,"protocol":"tcp","remote_group_id":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","remote_ip_prefix":"100.94.0.0/24","rule_zones":[{"resource_id":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","zone":"zone_name"}],"status":"created"}]'

## Other commands:
- Flags:
- --api-id string                 Api Id
- --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
- --error string                  Error
- --external-id string            External Id
- -h, --help                          help for update
- --project-type enum             project_type: Project type to update Security Group (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- --rules array(object)           Rules
- Use --rules=help for more details (default [])
- --security-group-id string      securityGroupId: Id of the async Security Group to update (required)
- --security-group-zones object   Security Group Zones
- Use --security-group-zones=help for more details (default {})
- --status string                 Status (required)
- -v, --version                       version for update

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri         Manually specify the server to use
      --x-request-id string    X-Request-Id
```

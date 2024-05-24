# This route is for internaly to update the information
when finished to request a instance creation on URP.

## Usage:
```bash
After requested successfully this route is called to save
the network information and urp instance ID.
```

## Product catalog:
- ### Note
- This route is used only for internal proposes.

## Other commands:
- Usage:
- ./mgc virtual-machine-xaas instances internal-instances create [id] [flags]

## Flags:
```bash
Flags:
      --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
      --error string                  Error
  -h, --help                          help for create
      --id uuid                       Id (required)
      --network-ids array(object)     Network Ids
                                      Use --network-ids=help for more details (default [])
      --project-type enum             Project Type (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (default "default")
      --status string                 Status (required)
      --urp-instance-id string        Urp Instance Id
  -v, --version                       version for create
```


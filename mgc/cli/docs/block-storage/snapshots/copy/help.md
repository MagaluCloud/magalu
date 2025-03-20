# Copy a object snapshot cross region for the currently authenticated tenant.

## Usage:
```bash
#### Rules
- The copy only be accepted when the destiny region is different from origin region.
- The copy only be accepted if the snapshot name in destiny region is different from input name.
- The copy only be accepted if the user has access to destiny region.
```

## Product catalog:
- #### Notes
- - Utilize the **block-storage snapshots list** command to retrieve a list of
- all Snapshots and obtain the ID of the Snapshot you wish to copy across different region.

## Other commands:
- Usage:
- mgc block-storage snapshots copy [id] [flags]

## Flags:
```bash
Flags:
      --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
      --destination-region enum       Regions (one of "br-mgl1", "br-ne1" or "br-se1") (required)
  -h, --help                          help for copy
      --id uuid                       Id (required)
  -v, --version                       version for copy
```


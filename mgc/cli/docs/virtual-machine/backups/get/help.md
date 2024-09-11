# Get a backup details for the current tenant which is logged in.

## Usage:
```bash
#### Notes
- You can use the backup list command to retrieve all backups,
so you can get the id of the backup that you want to get details.
```

## Product catalog:
- - You can use the **expand** argument to get more details from the object
- like instance.

## Other commands:
- Usage:
- ./mgc virtual-machine backups get [id] [flags]

## Flags:
```bash
Flags:
      --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
      --expand array(string)          Expand: You can get more detailed info about: ['instance']  (default [])
  -h, --help                          help for get
      --id string                     Id (required)
  -v, --version                       version for get
```


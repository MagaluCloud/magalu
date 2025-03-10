# Rename a Volume for the currently authenticated tenant.

## Usage:
```bash
#### Rules
- The Volume name must be unique; otherwise, renaming will not be allowed.
```

## Product catalog:
- #### Notes
- - Utilize the **block-storage volume list** command to retrieve a list of all
- Volumes and obtain the ID of the Volume you wish to rename.

## Other commands:
- Usage:
- ./mgc block-storage volumes rename [id] [flags]

## Flags:
```bash
Flags:
      --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
      --cli.watch                     Wait until the operation is completed by calling the 'get' link and waiting until termination. Akin to '! get -w'
  -h, --help                          help for rename
      --id uuid                       Id (required)
      --name string                   Name (between 3 and 50 characters) (required)
  -v, --version                       version for rename
```


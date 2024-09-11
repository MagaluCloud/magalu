# Patches a Backup for the currently authenticated tenant.

## Usage:
```bash
#### Rules
- The Backup name must be unique; otherwise, renaming will not be allowed.
- The Backup's state must be available.
```

## Product catalog:
- #### Notes
- - Utilize the **block-storage backups list** command to retrieve a list of
- all Backups and obtain the ID of the Backup you wish to rename.

## Other commands:
- Usage:
- ./mgc block-storage backups rename [id] [flags]

## Flags:
```bash
Flags:
      --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
      --description string            Description (between 1 and 255 characters)
  -h, --help                          help for rename
      --id uuid                       Id (required)
      --name string                   Name (between 1 and 255 characters)
  -v, --version                       version for rename
```


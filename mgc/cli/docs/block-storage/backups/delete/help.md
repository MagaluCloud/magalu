# Delete a Backup for the currently authenticated tenant.

## Usage:
```bash
#### Rules
- The Backup's status must be "completed".
- The Backup's state must be "available".
```

## Product catalog:
- #### Notes
- - Utilize the **block-storage backups** list command to retrieve a list of
- all Backups and obtain the ID of the Backup you wish to delete.

## Other commands:
- Usage:
- ./mgc block-storage backups delete [id] [flags]

## Flags:
```bash
Flags:
  -h, --help      help for delete
      --id uuid   Id (required)
  -v, --version   version for delete
```


# Create a backup for the currently authenticated tenant.

## Usage:
```bash
The Backup can be used when it reaches the "available" state and the
 "completed" status.
```

## Product catalog:
- #### Rules
- - The Backup name must be unique; otherwise, the creation will be disallowed.
- - The Volume can be either in in-use or available states.
- - The Volume must not have an operation in execution.

## Other commands:
- #### Notes
- - Use the **block-storage volume list** command to retrieve a list of all
- Volumes and obtain the ID of the Volume that will be used to create the
- Backup.

## Flags:
```bash
Usage:
  ./mgc block-storage backups create [flags]
```


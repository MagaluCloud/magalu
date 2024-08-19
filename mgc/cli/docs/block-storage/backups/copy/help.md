# Copy a backup cross region for the currently authenticated tenant.

## Usage:
```bash
#### Rules
- The copy only be accepted when the destiny region is different from origin region.
- The copy only be accepted if the backup's name in destiny region is different from input name.
- The copy only be accepted if the user has access to destiny region.
```

## Product catalog:
- #### Notes
- - Utilize the **block-storage backups list** command to retrieve a list of
- all Backups and obtain the ID of the Backup you wish to copy across different region.

## Other commands:
- Usage:
- ./mgc block-storage backups copy [flags]

## Flags:
```bash
Flags:
      --backup object           BackupIdRequest (properties: id and name)
                                Use --backup=help for more details (required)
      --backup.id string        BackupIdRequest: Id (between 1 and 255 characters)
                                This is the same as '--backup=id:string'.
      --backup.name string      BackupIdRequest: Name (between 1 and 255 characters)
                                This is the same as '--backup=name:string'.
      --destiny-region string   Destiny Region (between 1 and 255 characters) (required)
  -h, --help                    help for copy
  -v, --version                 version for copy
```


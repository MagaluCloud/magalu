# Deletes permanently an instance with the id provided in the current tenant
which is logged in.

## Usage:
```bash
#### Notes
- You can use the virtual-machine list command to retrieve all instances, so
- you can get the id of the instance that you want to delete.
```

## Product catalog:
- #### Result
- - The attached ports will be deleted as well.</li>
- - The attached volumes will be detached.</li>

## Other commands:
- Usage:
- ./mgc virtual-machine instances delete [id] [flags]

## Flags:
```bash
Flags:
      --delete-public-ip   Delete Public Ip
  -h, --help               help for delete
      --id uuid            Id (required)
  -v, --version            version for delete
```


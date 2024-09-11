# Create a backup of a Virtual Machine with the current tenant which is logged in.

## Usage:
```bash
A Backup is ready for restore when it's in completed status.
```

## Product catalog:
- #### Rules
- - It's possible to create a maximum of 100 backups per virtual machine.
- - In case quota reached, choose a backup to remove.
- - You can inform ID or Name from a Virtual Machine if both informed the priority will be ID.
- - It's only possible to create a backup of a valid virtual machine.
- - Each backup must have a unique name. It's not possible to create backups with the same name.

## Other commands:
- Usage:
- ./mgc virtual-machine backups create [flags]

## Flags:
```bash
Examples:
  ./mgc virtual-machine backups create --instance.id="9ec75090-2872-4f51-8111-53d05d96d2c6" --instance.name="some_resource_name"
```


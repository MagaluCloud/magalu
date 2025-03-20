# Restore a snapshot of an instance with the current tenant which is logged in. </br>

## Usage:
```bash
#### Notes
- You can check the snapshot state using snapshot list command.
- Use "machine-types list" to see all machine types available.
```

## Product catalog:
- #### Rules
- - A Snapshot is ready to restore when it's in available state.
- - To restore a snapshot you have to inform the new instance settings.
- - You must choose a machine-type that has a disk equal or larger
- than the original instance.

## Other commands:
- Usage:
- mgc virtual-machine snapshots restore [id] [flags]

## Flags:
```bash
Examples:
  mgc virtual-machine snapshots restore --machine-type.id="9ec75090-2872-4f51-8111-53d05d96d2c6" --machine-type.name="some_resource_name" --network.associate-public-ip=true --network.interface.id="9ec75090-2872-4f51-8111-53d05d96d2c6" --network.interface.security-groups='[{"id":"9ec75090-2872-4f51-8111-53d05d96d2c6"}]' --network.vpc.id="9ec75090-2872-4f51-8111-53d05d96d2c6" --network.vpc.name="some_resource_name"
```


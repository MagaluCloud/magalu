# WARN	magalu.cloud/sdk/openapi.virtual-machine.snapshots.create	ignored broken link	{"link": "delete", "error": "linked operationId=\"delete_instance_type_v0_instance_types__instance_type_id__delete\": could not resolve \"/operationIds/delete_instance_type_v0_instance_types__instance_type_id__delete\": missing field: \"delete_instance_type_v0_instance_types__instance_type_id__delete\""}
Create a snapshot of a Virtual Machine in the current tenant which is logged in. </br>
A Snapshot is ready for restore when it's in available state.

## Usage:
```bash
#### Notes
- You can verify the state of snapshot using the snapshot get command,
- To create a snapshot it's mandatory inform a valid and unique name.
```

## Product catalog:
- #### Rules
- - It's only possible to create a snapshot of a valid virtual machine.
- - It's not possible to create 2 snapshots with the same name.
- - You can inform ID or Name from a Virtual Machine if both informed the priority will be ID.

## Other commands:
- Usage:
- ./cli virtual-machine snapshots create [flags]

## Flags:
```bash
Examples:
  ./cli virtual-machine snapshots create --virtual-machine.id="9ec75090-2872-4f51-8111-53d05d96d2c6" --virtual-machine.name="some_resource_name"
```


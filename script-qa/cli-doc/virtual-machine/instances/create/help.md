# Creates a Virtual Machine instance in the current tenant which is logged in.

## Usage:
```bash
An instance is ready for you to use when it's in the running state.
```

## Product catalog:
- #### Notes
- - For the image data, you can use the virtual-machine images list command
- to list all available images.
- - For the machine type data, you can use the virtual-machine machine-types
- list command to list all available machine types.
- - You can verify the state of your instance using the virtual-machine get
- command.

## Other commands:
- #### Rules

## Flags:
```bash
- If you don't specify a VPC, the default VPC will be used. When the
default VPC is not available, the command will fail.
- If you don't specify an network interface, an default network interface
will be created.
- You can either specify an image id or an image name. If you specify
both, the image id will be used.
- You can either specify a machine type id or a machine type name. If
you specify both, the machine type id will be used.
- You can either specify an VPC id or an VPC name. If you specify both,
the VPC id will be used.
- The user data must be a Base64 encoded string.
```


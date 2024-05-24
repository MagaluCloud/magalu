# Register a URP Image on Virtual Machine DB.

## Usage:
```bash
### Note
The Image on URP need to be public and protected.
```

## Product catalog:
- Usage:
- ./mgc virtual-machine-xaas images create [flags]

## Other commands:
- Examples:
- ./mgc virtual-machine-xaas images create --end-life-at="2022-01-01T00:00:10Z" --end-standard-support-at="2022-01-01T00:00:10Z" --release-at="2022-01-01T00:00:10Z"

## Flags:
```bash
Flags:
      --cli.list-links enum[=table]      List all available links for this command (one of "json", "table" or "yaml")
      --end-life-at string               End Life At
      --end-standard-support-at string   End Standard Support At
  -h, --help                             help for create
      --image-id string                  Image Id (required)
      --image-url uri                    Image Url (between 1 and 2083 characters)
      --internal                         Internal (required)
      --min-disk integer                 Min Disk (min: 1) (required)
      --min-ram integer                  Min Ram (min: 1) (required)
      --min-vcpu integer                 Min Vcpu (min: 1) (required)
      --name string                      Name (between 1 and 255 characters) (required)
      --param.version string             Version
      --platform enum                    ImageV1Platform (one of "linux" or "windows")
      --release-at string                Release At
      --sku string                       Sku (required)
      --status enum                      (one of "active", "creating", "deleted", "deprecated", "error", "importing" or "pending") (default "active")
  -v, --version                          version for create
```


# This route is when the vm-instace worker already
deleted the instance from: urp and vpc api to mark the instance

## Usage:
```bash
to 'deleted' on virtual machine DB.
```

## Product catalog:
- ### Note
- This route is used only for internal proposes.

## Other commands:
- Usage:
- ./mgc virtual-machine-xaas instances internal-instances delete-id [id] [flags]

## Flags:
```bash
Flags:
  -h, --help                help for delete-id
      --id uuid             Id (required)
      --project-type enum   Project Type (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (default "default")
  -v, --version             version for delete-id
```


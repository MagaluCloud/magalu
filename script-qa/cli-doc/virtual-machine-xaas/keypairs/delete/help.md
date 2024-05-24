# Delete a existing keypair direct on urp using a informed tenant id.

## Usage:
```bash
### Note
This route is used only for internal proposes.
```

## Product catalog:
- Usage:
- ./mgc virtual-machine-xaas keypairs delete [keypair-name] [flags]

## Other commands:
- Flags:
- -h, --help                  help for delete
- --keypair-name string   Keypair Name (required)
- --project-type enum     ProjectTypeAll (one of "dbaas", "default", "iamaas", "k8saas" or "mngsvc") (required)
- -v, --version               version for delete

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri         Manually specify the server to use
```


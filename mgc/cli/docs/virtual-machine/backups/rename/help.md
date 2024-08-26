# Renames a backup with the id provided in the current tenant which is logged in.

## Usage:
```bash
#### Rules
- You can use the backup list command to retrieve all backups, so you can get the id of
the backup that you want to rename.
- A Backup can only renamed when it's in completed status.
- Each backup must have a unique name.
```

## Product catalog:
- Usage:
- ./mgc virtual-machine backups rename [id] [flags]

## Other commands:
- Flags:
- --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
- --cli.watch                     Wait until the operation is completed by calling the 'get' link and waiting until termination. Akin to '! get -w'
- -h, --help                          help for rename
- --id string                     Id (required)
- --name string                   Name (between 1 and 255 characters) (required)
- -v, --version                       version for rename

## Flags:
```bash
Global Flags:
      --api-key string           Use your API key to authenticate with the API
  -U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                                 use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                                 a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
  -t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                                 Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
      --debug                    Display detailed log information at the debug level
      --env enum                 Environment to use (one of "pre-prod" or "prod") (default "prod")
      --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
  -o, --output string            Change the output format. Use '--output=help' to know more details. (default "yaml")
  -r, --raw                      Output raw data, without any formatting or coloring
      --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri           Manually specify the server to use
```


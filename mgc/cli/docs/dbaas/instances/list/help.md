# Returns a list of database instances for a x-tenant-id.

## Usage:
```bash
Usage:
  ./mgc dbaas instances list [flags]
```

## Product catalog:
- Examples:
- ./mgc dbaas instances list --status="ACTIVE"

## Other commands:
- Flags:
- --control.expand enum       Instance extra attributes or relations to show with the main query. When available, more than one value can be informed using commas. e.g: '--control.expand="replicas"' (must be "replicas")
- --control.limit integer     The maximum number of items per page. (range: 1 - 25) (default 10)
- --control.offset integer    The number of items to skip before starting to collect the result set. (min: 0)
- --engine-id uuid            Engine Id unique identifier
- -h, --help                      help for list
- --status enum               Value referring to instance status. (one of "ACTIVE", "BACKING_UP", "CREATING", "DELETED", "DELETING", "ERROR", "ERROR_DELETING", "MAINTENANCE", "PENDING", "REBOOT", "RESIZING", "RESTORING", "STARTING", "STOPPED" or "STOPPING")
- -v, --version                   version for list
- --volume.size integer       Volume Size exact size
- --volume.size-gt integer    Volume Size greater than
- --volume.size-gte integer   Volume Size greater than or equal
- --volume.size-lt integer    Volume Size less than
- --volume.size-lte integer   Volume Size less than or equal

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
  -o, --output string            Change the output format. Use '--output=help' to know more details.
  -r, --raw                      Output raw data, without any formatting or coloring
      --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri           Manually specify the server to use
```


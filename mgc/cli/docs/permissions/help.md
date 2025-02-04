# CHANGE-ME

## Usage:
```bash
Usage:
  ./mgc permissions [flags]
  ./mgc permissions [command]
```

## Product catalog:
- Commands:
- health                    Operations related to API health checks.
- onboarding                Operations related to onboarding.
- permission-model-status   Operations related to view permission model status.
- permission-models         Operations related to permission models.
- permissions               Operations related to permissions within models.
- privileges-users          Operations related to privileges users.
- products-resources        Operations related to products resources list.
- resource-privileges       Operations related to resource privileges list.
- resource-privileges-check Operations related to resource privileges check.
- user-delete               Operations related to user delete.

## Other commands:
- Flags:
- -h, --help   help for permissions

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
      --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
  -o, --output string            Change the output format. Use '--output=help' to know more details.
  -r, --raw                      Output raw data, without any formatting or coloring
```


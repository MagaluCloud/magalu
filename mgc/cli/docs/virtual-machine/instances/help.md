# Operations with instances, including create, delete, start, stop, reboot and other actions.

## Usage:
```bash
mgc virtual-machine instances [flags]
  mgc virtual-machine instances [command]
```

## Commands:
```bash
create            Create an instance.
  delete            Delete an instance.
  get               Retrieve the details of a specific instance.
  list              List all instances.
  network-interface network-interface
  password          Retrieve the first windows admin password
  reboot            Reboot an instance.
  rename            Renames an instance.
  retype            Changes an instance machine-type.
  start             Starts an instance.
  stop              Stops an instance.
  suspend           Suspends instance.

```

## Flags:
```bash
  -h, --help      help for instances
  -v, --version   version for instances

```

## Global Flags:
```bash
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


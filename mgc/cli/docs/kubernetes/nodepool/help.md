# Endpoints related to the creation, listing of nodepools and nodes, updating, and deletion of nodepools for a Kubernetes cluster.

## Usage:
```bash
Usage:
  mgc kubernetes nodepool [flags]
  mgc kubernetes nodepool [command]
```

## Product catalog:
- Commands:
- create      Create a node pool
- delete      Delete node pool by node_pool_id
- get         Get node pool by node_pool_id
- list        List node pools by cluster_id
- nodes       List nodes from a node pool by node_pool_id
- update      Patch node pool replicas by node_pool_id

## Other commands:
- Flags:
- -h, --help      help for nodepool
- -v, --version   version for nodepool

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


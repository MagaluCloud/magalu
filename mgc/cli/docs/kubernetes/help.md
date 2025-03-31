# APIs related to the Kubernetes product.

## Usage:
```bash
mgc kubernetes [flags]
  mgc kubernetes [command]
```

## Commands:
```bash
cluster     Endpoints related to the creation, listing, deletion, and retrieval of the kubeconfig for a Kubernetes cluster.
  flavor      Endpoints related to listing available flavors for a Kubernetes cluster.
  info        Endpoints related to listing available flavors (resource configurations) and Kubernetes versions for a Kubernetes cluster.
  nodepool    Endpoints related to the creation, listing of nodepools and nodes, updating, and deletion of nodepools for a Kubernetes cluster.
  version     Endpoints related to listing available Kubernetes versions for a Kubernetes cluster.

```

## Flags:
```bash
  -h, --help      help for kubernetes
  -v, --version   version for kubernetes

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


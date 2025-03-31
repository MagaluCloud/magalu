# Endpoints related to the creation, listing, deletion, and retrieval of the kubeconfig for a Kubernetes cluster.

## Usage:
```bash
mgc kubernetes cluster [flags]
  mgc kubernetes cluster [command]
```

## Commands:
```bash
create      Create a cluster
  delete      Delete a cluster
  get         List a cluster by uuid
  kubeconfig  Get kubeconfig cluster
  list        List all clusters
  update      Patches a cluster

```

## Flags:
```bash
-h, --help      help for cluster
  -v, --version   version for cluster

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


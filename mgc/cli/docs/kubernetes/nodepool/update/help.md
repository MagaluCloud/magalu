# Updates nodes from a node pool by nodepool_uuid.

## Usage:
```bash
Usage:
  mgc kubernetes nodepool update [cluster-id] [node-pool-id] [flags]
```

## Product catalog:
- Examples:
- mgc kubernetes nodepool update --auto-scale.max-replicas=5 --auto-scale.min-replicas=2

## Other commands:
- Flags:
- --auto-scale object                 Object specifying properties for updating workload resources in the Kubernetes cluster.
- (properties: max_replicas and min_replicas)
- Use --auto-scale=help for more details
- --auto-scale.max-replicas integer   Object specifying properties for updating workload resources in the Kubernetes cluster: Maximum number of replicas for autoscaling. If not provided, the autoscale value will be assumed based on the "replicas" field.
- (min: 0)
- This is the same as '--auto-scale=max_replicas:integer'.
- --auto-scale.min-replicas integer   Object specifying properties for updating workload resources in the Kubernetes cluster: Minimum number of replicas for autoscaling. If not provided, the autoscale value will be assumed based on the "replicas" field.
- (min: 0)
- This is the same as '--auto-scale=min_replicas:integer'.
- --cli.list-links enum[=table]       List all available links for this command (one of "json", "table" or "yaml")
- --cluster-id uuid                   Cluster's UUID. (required)
- -h, --help                              help for update
- --node-pool-id uuid                 Nodepool's UUID. (required)
- --replicas integer                  Number of replicas of the nodes in the node pool.
- -v, --version                           version for update

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


# Creates a node pool in a Kubernetes cluster.

## Usage:
```bash
Usage:
  ./mgc kubernetes nodepool create [cluster-id] [flags]
```

## Product catalog:
- Examples:
- ./mgc kubernetes nodepool create --auto-scale.max-replicas=5 --auto-scale.min-replicas=2 --flavor="cloud-k8s.gp1.small" --name="nodepool-example" --replicas=3 --tags='["tag-value1"]'

## Other commands:
- Flags:
- --auto-scale object                 Object specifying properties for updating workload resources in the Kubernetes cluster.
- (properties: max_replicas and min_replicas)
- Use --auto-scale=help for more details
- --auto-scale.max-replicas integer   Object specifying properties for updating workload resources in the Kubernetes cluster: Maximum number of replicas for autoscaling. If not provided, or if the value is equal to or less than zero, the autoscaler value will be assumed based on the "replicas" field.
- (min: 1)
- This is the same as '--auto-scale=max_replicas:integer'. (default 1)
- --auto-scale.min-replicas integer   Object specifying properties for updating workload resources in the Kubernetes cluster: Minimum number of replicas for autoscaling. If not provided, or if the value is equal to or less than zero, the autoscaler value will be assumed based on the "replicas" field.
- (min: 1)
- This is the same as '--auto-scale=min_replicas:integer'. (default 1)
- --cli.list-links enum[=table]       List all available links for this command (one of "json", "table" or "yaml")
- --cluster-id uuid                   Cluster's UUID. (required)
- --flavor string                     Definition of the CPU, RAM, and storage capacity of the nodes.
- 
- | Flavor                     | vCPUs | RAM (GB) | Root Disk (GB) |
- |----------------------------|-------|----------|----------------|
- | cloud-k8s.gp1.small        | 2     | 4        | 20             |
- | cloud-k8s.gp1.medium       | 4     | 8        | 50             |
- | cloud-k8s.gp1.large        | 8     | 16       | 100            |
- | cloud-k8s.gp2.small        | 2     | 4        | 300            |
- | cloud-k8s.gp2.large        | 8     | 32       | 300            |
- | cloud-k8s.hm1.medium       | 4     | 16       | 50             |
- | cloud-k8s.i1-c48-r96-d300  | 48    | 96       | 300            |
- | cloud-k8s.i1-c1-r8-d100    | 1     | 8        | 100            |
- | cloud-k8s.i1-c8-r8-d40     | 8     | 8        | 40             |
- | cloud-k8s.i1-c4-r16-d100   | 4     | 16       | 100            |
- | cloud-k8s.i1-c4-r32-d100   | 4     | 32       | 100            |
- | cloud-k8s.i1-c8-r32-d100   | 8     | 32       | 100            |
- | cloud-k8s.i1-c16-r32-d40   | 16    | 32       | 40             |
- | cloud-k8s.i1-c16-r64-d100  | 16    | 64       | 100            |
- | cloud-k8s.i1-c32-r64-d100  | 32    | 64       | 100            |
- | cloud-k8s.i1-c32-r128-d500 | 32    | 128      | 500            |
- (required)
- -h, --help                              help for create
- --name string                       Name of the node pool. The name is primarily for idempotence and must be unique within a namespace. The name cannot be changed.
- The name must follow the following rules:
- - Must contain a maximum of 63 characters
- - Must contain only lowercase alphanumeric characters or '-'
- - Must start with an alphabetic character
- - Must end with an alphanumeric character
- (required)
- --replicas integer                  Number of replicas of the nodes in the node pool. (required) (default 1)
- --tags array(string)                List of tags applied to the node pool.
- --taints array(object)              Property associating a set of nodes.
- Use --taints=help for more details
- -v, --version                           version for create

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


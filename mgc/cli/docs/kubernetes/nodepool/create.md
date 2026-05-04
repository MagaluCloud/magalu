---
sidebar_position: 2
---
# Create

Creates a node pool in a Kubernetes cluster.

## Usage:
```
mgc kubernetes nodepool create [cluster-id] [flags]
```

## Examples:
```
mgc kubernetes nodepool create --auto-scale.max-replicas=5 --auto-scale.min-replicas=2 --flavor="BV2-4-40" --max-pods-per-node=32 --name="nodepool-example" --network.subnet-ids='["627c1f78-f9a6-4419-b582-a982144ff6bc","57b69486-e800-4cc3-92e2-6037b4cafe35"]' --replicas=3
```

## Flags:
```
    --auto-scale object                 Object specifying properties for updating workload resources in the Kubernetes cluster.
                                         (properties: max_replicas and min_replicas)
                                        Use --auto-scale=help for more details
    --auto-scale.max-replicas integer   Object specifying properties for updating workload resources in the Kubernetes cluster: Maximum number of replicas for autoscaling. If not provided, the autoscale value will be assumed based on the "replicas" field.
                                         (min: 0)
                                        This is the same as '--auto-scale=max_replicas:integer'.
    --auto-scale.min-replicas integer   Object specifying properties for updating workload resources in the Kubernetes cluster: Minimum number of replicas for autoscaling. If not provided, the autoscale value will be assumed based on the "replicas" field.
                                         (min: 0)
                                        This is the same as '--auto-scale=min_replicas:integer'.
    --availability-zones array(enum)    [Deprecated]List of availability zones where the resource can be created.
                                         
    --cli.list-links enum[=table]       List all available links for this command (one of "json", "table" or "yaml")
    --cluster-id uuid                   Cluster's UUID. (required)
    --flavor string                     Name of the machine type. The machine type defines the CPU, RAM, and storage capacity of the nodes.
                                        The full list of available machine types can be found in the Virtual Machine by listing the machine types.
                                        The smallest supported machine type is BV2-4-40.
                                        
                                        For V1 Clusters, the list of available flavors can be retrieved using the /v1/flavors endpoint (deprecated) or via the MGC CLI with the command 'kubernetes flavors list'.
                                         (required)
-h, --help                              help for create
    --max-pods-per-node integer         Maximum number of Pods allowed per node.
                                         (range: 8 - 110)
    --name string                       Name of the node pool. The name is primarily for idempotence and must be unique within a namespace. The name cannot be changed.
                                        The name must follow the following rules:
                                          - Must contain a maximum of 63 characters
                                          - Must contain only lowercase alphanumeric characters or '-'
                                          - Must start with an alphabetic character
                                          - Must end with an alphanumeric character
                                         (max character count: 63) (required)
    --network object                    Request object for the Kubernetes nodepools network resource request.
                                         (single property: subnet_ids)
                                        Use --network=help for more details
    --network.subnet-ids array          Request object for the Kubernetes nodepools network resource request: 
                                        This is the same as '--network=subnet_ids:array'.
    --replicas integer                  Number of replicas of the nodes in the node pool. (required) (default 1)
    --tags array(string)                [Deprecated]List of tags applied to the node pool. 
    --taints array(object)              Property associating a set of nodes.
                                        Use --taints=help for more details
```

## Global Flags:
```
    --api-key string           Use your API key to authenticate with the API
-U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                               use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                               a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
-t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                               Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
    --debug                    Display detailed log information at the debug level
    --env enum                 Environment to use (one of "pre-prod" or "prod") (default "prod")
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. You can use 'yaml', 'json' or 'table'.
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
```


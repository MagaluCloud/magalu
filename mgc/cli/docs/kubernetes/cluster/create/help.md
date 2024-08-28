# Creates a Kubernetes cluster in Magalu Cloud.

## Usage:
```bash
Usage:
  ./mgc kubernetes cluster create [flags]
```

## Product catalog:
- Examples:
- ./mgc kubernetes cluster create --allowed-cidrs='["192.168.1.0/24","10.0.0.0/16"]' --description="This is an example cluster." --enabled-bastion=false --enabled-server-group=false --name="cluster-example" --node-pools='[{"auto_scale":{"max_replicas":5,"min_replicas":2},"flavor":"cloud-k8s.gp1.small","name":"nodepool-example","replicas":3,"tags":["tag-value1"],"taints":[{"effect":"NoSchedule","key":"example-key","value":"valor1"}]}]' --param.version="v1.30.2" --zone="br-region-zone"

## Other commands:
- Flags:
- --allowed-cidrs array(string)   List of allowed CIDR blocks for API server access.
- 
- --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
- --description string            A brief description of the Kubernetes cluster.
- 
- --enabled-bastion               [Deprecated] This parameter is deprecated and its use won't create a bastion server
- Enables the use of a bastion host for secure access to the cluster.
- 
- --enabled-server-group          Enables the use of a server group with anti-affinity policy during the creation of the cluster and its node pools.
- (default true)
- -h, --help                          help for create
- --name string                   Kubernetes cluster name. The name is primarily intended for idempotence, and must be unique within a namespace. The name cannot be changed.
- The name must follow the following rules:
- - must contain a maximum of 63 characters
- - must contain only lowercase alphanumeric characters or '-'
- - must start with an alphabetic character
- - must end with an alphanumeric character
- (required)
- --node-pools array(object)      An array representing a set of nodes within a Kubernetes cluster.
- 
- Use --node-pools=help for more details
- --param.version string          The native Kubernetes version of the cluster.
- Please specify the Kubernetes version using the standard "vX.Y.Z" format.
- (default "v1.30.2")
- -v, --version                       version for create
- --zone string                   [Deprecated] This parameter is deprecated and its use won't create a cluster at requested zone.
- Identifier of the zone where the Kubernetes cluster will be located.

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


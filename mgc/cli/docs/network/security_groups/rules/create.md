# Create

Create a Rule async, returning its ID. To monitor the creation progress, please check the status in the service message or implement polling.Either a remote_ip_prefix or a remote_group_id can be specified.With remote_ip_prefix, all IPs that match the criteria will be allowed.With remote_group_id, only the specified security group is allowed to communicatefollowing the specified protocol, direction and port_range_min/max

## Usage:
```
mgc network security-groups rules create [security-group-id] [flags]
```

## Examples:
```
mgc network security-groups rules create --description="Allow incoming SSH traffic" --direction="ingress" --ethertype="IPv4"
```

## Flags:
```
    --cli.list-links enum[=table]   List all available links for this command (one of "json", "table" or "yaml")
    --description string            Description of the security group rule
    --direction string              Direction of the rule, either ingress or egress (required)
    --ethertype string              Ethertype of the rule, either IPv4 or IPv6 (required)
-h, --help                          help for create
    --port-range-max integer        Port Range Max
    --port-range-min integer        Port Range Min
    --protocol string               Protocol
    --remote-ip-prefix string       Remote Ip Prefix
    --security-group-id string      Security Group ID: Id of the Security Group (required)
    --validate-quota                validateQuota: Validate the quota before creating Rule
    --wait                          The request will be asynchronous. The wait parameter tells the API that you want the request to simulate synchronous behavior (to maintain endpoint compatibility). You can set an approximate timeout with the waitTimeout parameter
    --wait-timeout integer          waitTimeout: the approximate time in seconds you want to wait when simulating the request as synchronous (only works with wait=true)
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
-o, --output string            Change the output format. Use '--output=help' to know more details.
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
```


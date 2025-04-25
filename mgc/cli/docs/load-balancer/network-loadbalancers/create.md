# Create

Create Load Balancer

## Usage:
```
mgc load-balancer network-loadbalancers create [flags]
```

## Examples:
```
mgc load-balancer network-loadbalancers create --acls='[{"action":"ALLOW","ethertype":"IPv4","name":"acl for load balancer #1","protocol":"tcp","remote_ip_prefix":"192.168.67.0/24"}]' --backends='[{"balance_algorithm":"round_robin","description":"Some optional backend description 1","health_check_name":"nlb-health-check-1","name":"nlb-backend-1","targets":[{"nic_id":"00000000-0000-0000-0000-000000000000","port":80},{"nic_id":"00000000-0000-0000-0000-000000000001","port":443}],"targets_type":"instance"}]' --health-checks='[{"healthy_status_code":200,"name":"nlb-health-check-1","path":"/health-check","port":5000,"protocol":"tcp"}]' --listeners='[{"backend_name":"nlb-backend-1","name":"nlb-listener-1","port":80,"protocol":"tcp","tls_certificate_name":"nlb-tls-certificate-1"}]' --tls-certificates='[{"certificate":"SGVsbG8sIFdvcmxkIQ==","name":"nlb-tls-certificate-1","private_key":"SGVsbG8sIFdvcmxkIQ=="}]'
```

## Flags:
```
    --acls array(object)               Acls: The list of ACL configurations for the load balancer
                                       Use --acls=help for more details (default [])
    --backends array(object)           Backends: The list of backend configurations for the load balancer
                                       Use --backends=help for more details (required)
    --description string               A brief description of the load balancer
    --health-checks array(object)      Health Checks: The list of health check configurations for the load balancer
                                       Use --health-checks=help for more details (default [])
-h, --help                             help for create
    --listeners array(object)          Listeners: The list of listener configurations for the load balancer
                                       Use --listeners=help for more details (required)
    --name string                      The unique name of the load balancer (max character count: 64) (required)
    --panic-threshold integer          Panic Threshold: Minimum percentage of failed upstreams that load balancer will consider to give an alert (range: 0 - 100)
    --public-ip-id string              The public IP ID associated with the load balancer
    --subnet-pool-id string            The subnet pool ID associated with the load balancer
    --tls-certificates array(object)   The list of TLS certificates for the load balancer
                                       Use --tls-certificates=help for more details (default [])
    --type string                      The type of the load balancer (e.g., proxy) (default "proxy")
    --visibility enum                  The visibility of the load balancer (e.g., internal, external) (one of "external" or "internal") (required)
    --vpc-id string                    The VPC ID associated with the load balancer (required)
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


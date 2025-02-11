# Update Load Balancer by ID

## Usage:
```bash
Usage:
  ./mgc load-balancer application-loadbalancers replace [load-balancer-id] [flags]
```

## Product catalog:
- Examples:
- ./mgc load-balancer application-loadbalancers replace --acls='[{"action":"ALLOW","ethertype":"IPv4","id":"00000000-0000-0000-0000-000000000000","name":"acl for load balancer #1","port_range_max":81,"port_range_min":80,"protocol":"http","remote_ip_prefix":"192.168.67.10"}]' --backends='[{"balance_algorithm":"round_robin","description":"Some optional backend description 1","health_check_id":"00000000-0000-0000-0000-000000000000","health_check_name":"alb-health-check-1","id":"00000000-0000-0000-0000-000000000000","name":"alb-backend-1","targets":[{"id":"00000000-0000-0000-0000-000000000000","ip_address":"10.0.0.1","port":80},{"id":"00000000-0000-0000-0000-000000000001","ip_address":"10.0.0.2","port":443}],"targets_type":"raw"}]' --health-checks='[{"description":"Some optional health-check description 1","healthy_status_code":200,"healthy_threshold_count":8,"id":"00000000-0000-0000-0000-000000000000","initial_delay_seconds":30,"interval_seconds":30,"name":"alb-health-check-1","path":"/health-check","port":5000,"protocol":"http","timeout_seconds":10,"unhealthy_threshold_count":3}]' --listeners='[{"backend_id":"00000000-0000-0000-0000-000000000000","backend_name":"alb-backend-1","description":"Some optional listener description 1","id":"00000000-0000-0000-0000-000000000000","name":"alb-listener-1","port":443,"protocol":"https","public_ip_id":"00000000-0000-0000-0000-000000000000","tls_certificate_name":"alb-tls-certificate"}]' --tls-certificates='[{"certificate":"-----BEGIN CERTIFICATE-----","description":"Some optional tls-certificate description 1","id":"00000000-0000-0000-0000-000000000000","name":"alb-tls-certificate-1","private_key":"-----BEGIN PRIVATE KEY-----"}]'

## Other commands:
- Flags:
- --acls array(object)                  Acls
- Use --acls=help for more details (required)
- --backends array(object)              Backends
- Use --backends=help for more details (required)
- --cli.list-links enum[=table]         List all available links for this command (one of "json", "table" or "yaml")
- --description string                  Description
- --force-authentication                Force authentication by sending the header even if this API doesn't require it
- --health-checks array(object)         Health Checks
- Use --health-checks=help for more details (required)
- -h, --help                                help for replace
- --listeners array(object)             Listeners
- Use --listeners=help for more details (required)
- --load-balancer-id string             load_balancer_id: ID of the load balancer to update (required)
- --name string                         Name (required)
- --security-groups-ids array(string)   Security Groups Ids
- --subnet-pool-id string               Subnet Pool Id
- --tls-certificates array(object)      Tls Certificates
- Use --tls-certificates=help for more details (required)
- --type enum                           ApplicationLoadBalancerType (must be "application") (required)
- -v, --version                             version for replace
- --visibility enum                     LoadBalancerVisibility (one of "external" or "internal") (required)
- --vpc-id string                       Vpc Id (required)

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
      --server-url uri           Manually specify the server to use
```

